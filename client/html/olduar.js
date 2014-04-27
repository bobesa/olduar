var Server = "",
	Username = "",
	Password = "";

var FADE_ANIMATION_TIME = 500;

var SlotTypes = {
	"consumable": "Consumable",
	"1hand": "One handed",
	"2hand": "Two handed",
	"head": "Headwear",
	"torso": "Torso",
	"hands": "Handwear",
	"legs": "Leggings",
	"feet": "Footwear"
};

var request = function(method,command,value,callback){
	if(typeof value == "function") {
		callback = value;
		value = null;
	}
	var xhr = new XMLHttpRequest();
	xhr.addEventListener("readystatechange",function(){
		if(xhr.readyState==4) {
			if (xhr.status == 200) {
				document.getElementById("viewGame").style.display = "";
				document.getElementById("viewLogin").style.display = "none";
				if(callback) callback(JSON.parse(xhr.responseText));
			} else {
				document.getElementById("viewGame").style.display = "none";
				document.getElementById("viewLogin").style.display = "";
				if(callback) callback(null);
			}
		}
	});
	xhr.open(method.toUpperCase(), "/api/"+command);
	xhr.setRequestHeader("Authorization", "Basic " + btoa(Username+":"+Password));
	xhr.send(value);
},
	domEventLog,
	domLocationName,
	domLocationDescription,
	domLocationExits,
	domLocationActions,
	domLocationItems,
	listExits = {},
	listActions = {},
	listItems = {},
	listNpcs = {};

Array.prototype.hasOwnItemInPath = function(match,path){
	var len = this.length;
	for(var i=0;i<len;i++) {
		if(this[i][path] == match) return true;
	}
	return false;
};

var updateDomList = function(rootDom, list, data, update) {

	if(typeof data == "undefined" || data == null) {

		Object.keys(list).forEach(function(command){
			var entry = list[command];
			entry.dom.classList.add("fade");
			setTimeout(function(){
				rootDom.removeChild(entry.dom);
			},FADE_ANIMATION_TIME);
			delete list[command];
		});

	} else if(data instanceof Array) {

		//Remove unused entries
		Object.keys(list).forEach(function(command){
			if(!data.hasOwnItemInPath(command,"id")) {
				var entry = list[command];
				entry.dom.classList.add("fade");
				setTimeout(function(){
					rootDom.removeChild(entry.dom);
				},FADE_ANIMATION_TIME);
				delete list[command];
			}
		});

		//Update/Add entries from data
		data.forEach(function(item){
			//Prepare variables
			var command = item.id, description = item.name, entry = list.hasOwnProperty(command) ? list[command] : {};

			//Setup
			entry.id = command;
			entry.desc = description;
			if(!entry.dom) {
				entry.dom = document.createElement("div");
				rootDom.appendChild(entry.dom);
				update(entry);
			}
			entry.dom.innerHTML = command + " ("+description+")";

			//Append/Replace exit
			list[command] = entry;
		});

	} else if (data instanceof Object && data != null) {

		//Remove unused entries
		Object.keys(list).forEach(function(command){
			if(!data.hasOwnProperty(command)) {
				var entry = list[command];
				entry.dom.classList.add("fade");
				setTimeout(function(){
					rootDom.removeChild(entry.dom);
				},FADE_ANIMATION_TIME);
				delete list[command];
			}
		});

		//Update/Add entries from data
		Object.keys(data).forEach(function(command){
			//Prepare variables
			var description = data[command], entry = list.hasOwnProperty(command) ? list[command] : {};

			//Setup
			entry.id = command;
			entry.desc = description;
			if(!entry.dom) {
				entry.dom = document.createElement("div");
				rootDom.appendChild(entry.dom);
				update(entry);
			}
			entry.dom.innerHTML = command + " ("+description+")";

			//Append/Replace exit
			list[command] = entry;
		});

	}

};

var parseLocationData = function(data){
	if(data != null) {
		//Process location data
		domLocationName.innerHTML = data.name;
		domLocationDescription.innerHTML = data.desc;
		if(data.history) data.history.forEach(function(event){
			var dom = document.createElement("div");
			dom.innerHTML = event.text;
			domEventLog.appendChild(dom);
			domEventLog.scrollTop = domEventLog.scrollHeight;
		});

		updateDomList(domLocationExits,listExits,data.exits,function(entry){
			entry.dom.addEventListener("click",function(){
				request("get","go/"+entry.id,parseLocationData);
			});
		});
		updateDomList(domLocationActions,listActions,data.actions,function(entry){
			entry.dom.addEventListener("click",function(){
				request("get","do/"+entry.id,parseLocationData);
			});
		});
		updateDomList(domLocationItems,listItems,data.items,function(entry){
			entry.dom.addEventListener("click",function(){
				request("get","pickup/"+entry.id,parseLocationData);
			});
		});
		updateDomList(document.getElementById("locationNpcs"),listNpcs,data.npcs,function(entry){
			entry.dom.addEventListener("click",function(){
				request("post","attack/"+entry.id,parseLocationData);
			});
		});
	} else {
		//Show room selection
		toggleTabs("tabRooms");
	}

}, updateCredentials = function(){
	Username = document.getElementById("txtUser").value;
	Password = document.getElementById("txtPass").value;
	window.localStorage.setItem("olduar_username",Username);
	window.localStorage.setItem("olduar_password",Password);

}, getQualityColor = function(quality) {
	return "#000";

},toggleTabs = function(wanted){
	//Show wanted tab - hide others
	["tabLocation","tabInventory","tabRooms","tabAccount"].forEach(function(name){ document.getElementById(name).style.display = name==wanted?"block":"none"; });
	if(wanted == "tabInventory") {
		request("GET","inventory",function(inventory){
			var domInv = document.getElementById("inventory");
			if(inventory == null || inventory.length == 0) {
				domInv.innerHTML = "Your inventory is empty";
			} else {
				domInv.innerHTML = "";
				inventory.forEach(function(item){
					var dom = document.createElement("div");
					dom.className = "item";
					dom.innerHTML = '<span class="name" style="color:'+getQualityColor(item.quality)+'">' + item.name + "</span><br>" + '<span class="desc">'+item.desc+'</span>';
					if(item.usable) dom.addEventListener("click",function(){
						request("post","use/"+item.id,parseLocationData);
					});
					if(item.type != "consumable" && item.type != "") dom.addEventListener("click",function(){
						request("post","equip/"+item.id,parseLocationData);
					});
					domInv.appendChild(dom);

					//Attach events
					var domInspect = document.getElementById("inspect");
					dom.addEventListener("mouseover",function(e){
						domInspect.style.left = e.clientX+"px";
						domInspect.style.top = (e.clientY+20)+"px";
						domInspect.style.display = "block";
						domInspect.style.innerHTML = "Loading...";
						request("GET","inspect/"+item.id,function(detail){
							var str = '<div class="name" style="color:'+getQualityColor(item.quality)+'">'+detail.name+'</div>'; //TODO: Quality affect name color
							if(detail.type != "") str += '<div class="type">'+SlotTypes[detail.type]+(detail.class != "" && (detail.type == "consumable" || detail.type == "1hand" || detail.type == "2hand") ? " " + detail.class : "")+'</div>';
							if(detail.stats) Object.keys(detail.stats).forEach(function(stat){
								str += '<div class="stat"><b>'+stat+'</b>: ';
								var values = detail.stats[stat];
								if(values.min == values.max) {
									str += values.min;
								} else {
									str += values.min + " - " + values.max;
								}
								str += '</div>';
							});
							if(detail.desc != "") str += '<div class="desc">"'+detail.desc+'"</div>';
							if(detail.weight > 0) str += '<div class="detail">Weights '+detail.weight+'lb</div>';
							if(detail.usable) str += '<div class="detail">This is usable item</div>';
							domInspect.innerHTML = str;
						});
					});
					dom.addEventListener("mouseout",function(){
						domInspect.style.display = "none";
					});
					dom.addEventListener("mousemove",function(){
						domInspect.style.left = event.clientX+"px";
						domInspect.style.top = (event.clientY+20)+"px";
					});
				});
			}
		});
	}

}, reloadListOfRooms = function(rooms){
	toggleTabs("tabRooms");
	var domList = document.getElementById("listRooms");
	domList.innerHTML = "";
	if(rooms.length == 0) {
		var dom = document.createElement("option");
		dom.innerHTML = "No rooms are currently active";
		dom.disabled = true;
		domList.appendChild(dom);

	} else {
		rooms.forEach(function(room){
			var dom = document.createElement("option");
			dom.innerHTML = room;
			dom.value = room;
			domList.appendChild(dom);
		});
	}
};

window.addEventListener("load",function(){
	domEventLog = document.getElementById("eventLog");
	domLocationName = document.getElementById("locationName");
	domLocationDescription = document.getElementById("locationDescription");
	domLocationExits = document.getElementById("locationExits");
	domLocationActions = document.getElementById("locationActions");
	domLocationItems = document.getElementById("locationItems");

	//List
	document.getElementById("listRooms").addEventListener("change",function(e){
		document.getElementById("txtRoom").value = this.value;
	});

	//Tabs
	document.getElementById("buttonTabLocation").addEventListener("click",function(){
		toggleTabs("tabLocation");
	});
	document.getElementById("buttonTabInventory").addEventListener("click",function(){
		toggleTabs("tabInventory");
	});
	document.getElementById("buttonTabRooms").addEventListener("click",function(){
		toggleTabs("tabRooms");
	});
	document.getElementById("buttonTabAccount").addEventListener("click",function(){
		toggleTabs("tabAccount");
	});

	//Buttons
	document.getElementById("buttonLogin").addEventListener("click",function(){
		updateCredentials();
		request("get","look",parseLocationData);
	});
	document.getElementById("buttonRegister").addEventListener("click",function(){
		updateCredentials();
		request("post","register",function(success){
			if(success){
				request("get","look",parseLocationData);
			} else {
				document.getElementById("viewGame").style.display = "none";
				document.getElementById("viewLogin").style.display = "";
				alert("Username is taken or invalid");
			}
		});
	});
	document.getElementById("buttonLogout").addEventListener("click",function(){
		request("post","leave",function(){
			Username = "";
			Password = "";
			document.getElementById("txtUser").value = Username;
			document.getElementById("txtPass").value = Password;
			document.getElementById("viewGame").style.display = "none";
			document.getElementById("viewLogin").style.display = "";
			window.localStorage.setItem("olduar_username",Username);
			window.localStorage.setItem("olduar_password",Password);
		});
	});
	document.getElementById("buttonRename").addEventListener("click",function(){
		request("post","rename",document.getElementById("txtName").value,function(success){
			if(!success) alert("Renaming failed")
		});
	});
	document.getElementById("buttonJoin").addEventListener("click",function(){
		toggleTabs("tabLocation");
		request("post","join/"+document.getElementById("txtRoom").value,parseLocationData);
	});
	document.getElementById("buttonLeave").addEventListener("click",function(){
		request("post","leave",function(rooms){
			reloadListOfRooms(rooms);
			//TODO: Handle option that player is forbidden to leave the room
		});
	});
	document.getElementById("buttonRefreshRooms").addEventListener("click",function(){
		request("get","rooms",reloadListOfRooms);
	});

	Username = window.localStorage.getItem("olduar_username") || "";
	Password = window.localStorage.getItem("olduar_password") || "";
	if(Username != "" && Password != "") {
		request("get","look",parseLocationData);
	}
});
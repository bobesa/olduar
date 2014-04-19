var Server = "",
	Username = "",
	Password = "";

var FADE_ANIMATION_TIME = 500;

var request = function(method,command,value,callback){
	if(typeof value == "function") {
		callback = value;
		value = null;
	}
	var xhr = new XMLHttpRequest();
	xhr.addEventListener("readystatechange",function(){
		if(xhr.readyState==4) {
			if (xhr.status == 200) {
				callback(JSON.parse(xhr.responseText));
				document.getElementById("game").style.display = "";
				document.getElementById("account").style.display = "none";
			} else {
				callback(null);
				document.getElementById("game").style.display = "none";
				document.getElementById("account").style.display = "";
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
	listItems = {};

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
			//TODO: Inspection on hover
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

	//Account settings
	Username = document.getElementById("txtUser").value;
	Password = document.getElementById("txtPass").value;

	//Buttons
	document.getElementById("buttonRegister").addEventListener("click",function(){
		request("post","register",function(success){
			if(success){
				request("get","look",parseLocationData);
			} else {
				alert("Username is taken or invalid")
			}
		});
	});
	document.getElementById("buttonRename").addEventListener("click",function(){
		request("post","rename",document.getElementById("txtName").value,function(success){
			if(!success) alert("Renaming failed")
		});
	});
	document.getElementById("buttonJoin").addEventListener("click",function(){
		request("post","join/"+document.getElementById("txtRoom").value,parseLocationData);
	});
	document.getElementById("buttonLeave").addEventListener("click",function(){
		request("post","leave",function(rooms){

		});
	});

	//Polling
	setInterval(function(){
		//request("get","look",parseLocationData);
	},2000);

	//Initial request
	request("get","look",parseLocationData);
});
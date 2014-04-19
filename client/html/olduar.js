var Server = "",
	Username = "test",
	Password = "test2";

var FADE_ANIMATION_TIME = 500;

var request = function(command,param,callback){
	if(param != "") command += "/";
	var xhr = new XMLHttpRequest();
	xhr.addEventListener("readystatechange",function(){
		if(xhr.readyState==4) {
			if (xhr.status == 200) {
				callback(JSON.parse(xhr.responseText));
			} else {
				callback(null);
			}
		}
	});
	xhr.open("GET", "/test/"+command+param);
	xhr.setRequestHeader("Authorization", "Basic " + btoa(Username+":"+Password));
	xhr.send();
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
				request("go",entry.id,parseLocationData);
			});
		});
		updateDomList(domLocationActions,listActions,data.actions,function(entry){
			entry.dom.addEventListener("click",function(){
				request("do",entry.id,parseLocationData);
			});
		});
		updateDomList(domLocationItems,listItems,data.items,function(entry){
			entry.dom.addEventListener("click",function(){
				request("pickup",entry.id,parseLocationData);
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
	setInterval(function(){
		//request("look","",parseLocationData);
	},2000);
	request("look","",parseLocationData);
});
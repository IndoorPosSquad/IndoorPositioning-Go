//alert("test");


window.onload=function(){  
	//Initialise Graph  
	var g = new canvasGraph('graph');  
	g.drawGraph([{x:0,y:200,z:0},
				 {x:1000,y:200,z:0}]);
    
};  

var refresh = function(x_d, y_d){
	var g = new canvasGraph('graph');  
	g.drawGraph([{x:0,y:200,z:0},
				 {x:1000,y:200,z:0}]);
	g.drawGraph([{x:x_d * 200, y: 200, z:y_d * 200}]);
};

try {
	var sock = new WebSocket("ws:/localhost:2000/sock");
//sock.binaryType = 'blob'; // can set it to 'blob' or 'arraybuffer
	console.log("Websocket - status: " + sock.readyState);
	sock.onopen = function(m) {
		console.log("CONNECTION opened..." + this.readyState);
	};
	sock.onmessage = function(m) {
		console.log("Incoming Msg:", m.data);
		var pos = m.data.split(',');
		console.log(pos);
		var x_p = parseFloat(pos[0]);
		x_p /= 100.0;
		var y_p = parseFloat(pos[1]);
		y_p /= 100.0;
		console.log(x_p, y_p);
		refresh(x_p, y_p);
	};
	sock.onerror = function(m) {
		console.log("Error occured sending..." + m.data);
	};
	sock.onclose = function(m) {
		console.log("Disconnected - status " + this.readyState);
	};
} catch(exception) {
	console.log(exception);
}

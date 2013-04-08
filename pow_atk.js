/**
 * An implemention of the proof-of-work protocol with
 * a puzzle called reverse compute hash.
 * This code solves the puzzle Reverse compute hash.
 * 
 * The code is in the public domain.
 * Written by: Felix Leopoldo Rios and Christian Hellman
 */

$(document).ready(function(){
    /* 
     * This search function is protected by proof-of-work.
     */
    table = $('#attackers');
	$.ajax({
	  url: "ip",
	  context: document.body,
	  success: function(response) {
	  				// alert(response);
				  $("#ip")[0].innerHTML = "[" + response + "]";
				}
	});
    $("#attackbtn").click(function(){
    	var d1 = $('#f1').data("delay_param")();
    	var d2 = $('#f2').data("delay_param")();		
		var behaviour_type = $('input[name=behaviour_type]').filter(':checked').val();

        var numWorkers = parseInt($('#number_of_attacker_field').val());
        startWorkerSwarm(numWorkers, d1, d2, behaviour_type);
    });
});

var table;
var startid = 0;

function getTD(id){
	var td =  $(document.createElement('td'));
	td.attr({'class': id});
	return td;
}

function buildRow(row){
	// row.append("<td>", {id: "local_id"});
	// row.append("<td>", {id: "remote_id"});
	// row.append("<td>", {id: "difficulty"});
	// row.append("<td>", {id: "number"});
	// row.append("<td>", {id: "abort"});
	row.append(getTD("local_id"));
	row.append(getTD("remote_id"));
	row.append(getTD("behaviour"));
	row.append(getTD("difficulty"));
	row.append(getTD("status"));
	row.append(getTD("solved"));
	row.append(getTD("requesting_ms")); // time for server to respond to request
	row.append(getTD("solving_ms"));	 // time for worker/client to solve a puzzle
	row.append(getTD("granting_ms"));	 // time for server to verify and reply
	row.append(getTD("close"));
	
}
function get(who, what){
	return who.children("."+what)[0];
}

function set(label, value){
	get(this, label).innerHTML = "" + value;
}

function rnd_snd() {
	return (Math.random()*2-1)+(Math.random()*2-1)+(Math.random()*2-1);
}

function delay(delay_param){
	//alert(val1);
	//
	var dist_type = delay_param.dtype, val1 = delay_param[0], val2 = delay_param[1]; 
	// log("val1: "+ val1 + " val2: " + val2 + "dist type: " + dist_type);
	var delay = 0;
	if(dist_type == "dist_stoch"){
		var rand = Math.round(rnd_snd()*val2+val1)
		if(rand < 0){
			delay = 0;
		}else{
			delay = Math.round(rnd_snd()*val2+val1);
		}

	}else if (dist_type == "dist_uni"){
		delay = Math.floor(Math.random() * (val2 - val1 + 1)) + val1;
	}
	return delay;
}


/*************************** Web worker logics ********************************/
/* 
 * Start webworkers
 */
function startWorkerSwarm(numWorkers, delay1, delay2, behaviour_type){
	// var sockets = [];
	// Initial setup.
		log("Starting " + numWorkers + " workers.");

		var dataminer = new WebSocket("ws://{{$}}/dataminer");

		for(var i = 0; i < numWorkers; i++){
			(function() {
				var requestingTime;
				var solvingTime;
				var grantingTime;
				var startRequestingTime;
				var startSolvingTime;
				var startGrantingTime;

				if (window["WebSocket"]) {
					var conn = new WebSocket("ws://{{$}}/ws");
					var w = new Worker("attacktask.js");
					var id = startid + i;
					var closed = false;
					var solved = 0;
					// log("Worker " + id + ": started on new websocket.");
					var trow = $(document.createElement('tr'));
					trow.set = set;
					table.append(trow);
					buildRow(trow);
					// trow.children("#local_id")[0].innerHTML = "" + id;
					trow.set("local_id", id);
					trow.set("behaviour", behaviour_type);
					trow.set("status", "CONNECTING");

					var b = $(document.createElement('button'));
					b.html("Close");
					var cl = $(get(trow, "close"));
					b.click(function(){
						conn.close();
					});
					cl.append(b);

					var response;
					// sockets.push(conn);

					//Setup Worker
					w.onmessage=function (e){
						if(closed) return;
						
						var endSolvingTime = Number(new Date().getTime());
						solvingTime = endSolvingTime - startSolvingTime;
						trow.set("solving_ms", solvingTime);
						
						trow.set("status", "SOLVED");
						solved++;
						trow.set("solved", solved);
							 // var worker_data=e.data;
						var solution = e.data.solution; 
							// recieved data from worker
							// alert(worker_data);
							//alert(JSON.stringify(solution));
						
		                // $("#result").append("<br/>" + solution + "<br/>");
		                var request = { "Problems": solution, 
		                                "Query": response.Query,
		                                "Hash": CryptoJS.SHA256(solution + "" + response["Seed"]).toString(CryptoJS.enc.Hex),
		                                "Opcode": 1};

		                //alert(JSON.stringify(request))

		                trow.set("status", "WAIT COMMIT")
		                setTimeout(function(){
		                	startGrantingTime = Number(new Date().getTime());
		                	conn.send(JSON.stringify(request));
		                },delay(delay2) );

			        } 

			        conn.onclose = function(evt) {  
			           // log("Connection closed.");
			           closed = true;
			           trow.set("status", "DISCONNECTED");
			        };
			        conn.onerror = function(evt) {  
			           // log("Connection closed.");
			           closed = true;
			           trow.set("status", "ERROR");
			        };
			        conn.onmessage = function(evt) {
			            // alert("Got response: " + §evt.data);

			            response = JSON.parse(evt.data);
			            if(response.SocketId ){
			            	trow.set("remote_id", response.SocketId);
			            }
			            
			            if(response["Opcode"] == 1){
			            	var endRequestingTime = Number(new Date().getTime());
							requestingTime = endRequestingTime - startRequestingTime;
							trow.set("requesting_ms", requestingTime);

			                // alert("Problems is:" + response.Problems);
			                // trow.children("#difficulty")[0].innerHTML = "" + response["Difficulty"];
			                trow.set("difficulty", "{" + response.Difficulty.Zeroes + " " + response.Difficulty.Problems + "}");

			                // trow.children("#number")[0].innerHTML = "" + response["Problems"].length;
			           		//Send message with data to worker
			           		trow.set("status", "WORKING");

			           		startSolvingTime = Number(new Date().getTime());
			           		w.postMessage({problems: response["Problems"], difficulty: response["Difficulty"].Zeroes, wait: delay(delay2)});
			            } else {
			            	var endGrantingTime = Number(new Date().getTime());
							grantingTime = endGrantingTime - startGrantingTime;
							trow.set("granting_ms", grantingTime);
			            	trow.set("status", "WAIT NEW");

			                var datacollection = { "Behaviour": behaviour_type, 
					                               "RequestingTime": requestingTime,
					                               "SolvingTime": solvingTime,
					                               "GrantingTime": grantingTime};

			            	setTimeout(function(){
								conn.onopen();
			                },delay(delay1) );
			           			
			           	}

			        }
			        conn.onopen = function(evt) {			        	
			        	trow.set("status", "CONNECTED");
			        	var request = {"Opcode": 0, "Query": "Calle"};

			        	startRequestingTime = Number(new Date().getTime());
			        	this.send(JSON.stringify(request));
			        }
			  	}
		    })();

		}
		startid+=numWorkers;
}


function log(msg){
	 $("#msg").append("<br/><b>"+msg+"</b><br/>");
}



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
        var numWorkers = parseInt($('#number_of_attacker_field').val());
        startWorkerSwarm(numWorkers)
    });
});

var table;
var startid = 0;

function getTD(id){
	var td =  $(document.createElement('td'));
	td.attr({'id': id});
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
	row.append(getTD("difficulty"));
	row.append(getTD("number"));
	row.append(getTD("status"));
}
function set(label, value){
	this.children("#"+label)[0].innerHTML = "" + value;
}

/*************************** Web worker logics ********************************/
/* 
 * Start webworkers
 */
function startWorkerSwarm(numWorkers){
	// var sockets = [];
	// Initial setup.
		log("Starting " + numWorkers + " workers.");


		for(var i = 0; i < numWorkers; i++){
			(function() {
				if (window["WebSocket"]) {
					var conn = new WebSocket("ws://{{$}}/ws");
					var w = new Worker("attacktask.js");
					var id = startid + i;
					log("Worker " + id + ": started on new websocket.");
					var trow = $(document.createElement('tr'));
					trow.set = set;
					table.append(trow);
					buildRow(trow);
					// trow.children("#local_id")[0].innerHTML = "" + id;
					trow.set("local_id", id);


					var response;
					// sockets.push(conn);

					//Setup Worker
					w.onmessage=function (e){
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
		                conn.send(JSON.stringify(request));
			        } 

			        conn.onclose = function(evt) {  
			           log("Connection closed.");
			        };
			        conn.onmessage = function(evt) {
			            // alert("Got response: " + §evt.data);
			            response = JSON.parse(evt.data);
			            
			            if(response["Opcode"] == 1){

			                // alert("Problems is:" + response.Problems);
			                // trow.children("#difficulty")[0].innerHTML = "" + response["Difficulty"];
			                trow.set("difficulty", response["Difficulty"]);
			                trow.set("number", response["Problems"].length);

			                // trow.children("#number")[0].innerHTML = "" + response["Problems"].length;
			           		//Send message with data to worker
			           		w.postMessage({problems: response["Problems"], difficulty: response["Difficulty"]});
			            } else {
			           			conn.onopen();
			           	}

			        }
			        conn.onopen = function(evt) {
			        	var request = {"Opcode": 0, "Query": "Calle"};
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



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
    $("#attackbtn").click(function(){
        var numWorkers = parseInt($('#number_of_attacker_field').val());
        startWorkerSwarm(numWorkers)
    });
});


/*************************** Web worker logics ********************************/
/* 
 * Start webworkers
 */
function startWorkerSwarm(numWorkers){
	// var sockets = [];
	// Initial setup.
		for(var i = 0; i < numWorkers; i++){
			(function() {
				if (window["WebSocket"]) {
					var conn = new WebSocket("ws://{{$}}/ws");
					var w = new Worker("attacktask.js");
					var response;
					// sockets.push(conn);

					//Setup Worker
					w.onmessage=function (e){
								 var worker_data=e.data;
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
}


function log(msg){
	 $("#msg").append("<br/><b>"+msg+"</b><br/>");
}
function append_result(str){
	$("#result").append(str);
}



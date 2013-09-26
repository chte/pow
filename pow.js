/**
 * An implemention of the proof-of-work protocol with
 * a puzzle called reverse compute hash.
 * This code solves the puzzle Reverse compute hash.
 * 
 * The code is in the public domain.
 * Written by: Felix Leopoldo Rios and Christian Hellman
 */


function randomString(sChrs,iLen) {
    var sRnd = '';
    for (var i=0; i<iLen;i++){
	var randomPoz = Math.floor(Math.random() * sChrs.length);
	sRnd += sChrs.substring(randomPoz,randomPoz+1);
    }
    return sRnd;
}


/** 
 * This function actually computes the puzzle, i.e.
 * tries to find an x which makes h(x,seed) to a string with 
 * the specified difficulty leading zeros.
 */
function find_x(difficulty, seed){
    var cmp='';
    for(var i=0;i<difficulty;i++){
	cmp+='0';
    }
    var x=0;
    while (true){
	str = x + "" + seed;
	if(CryptoJS.SHA256(str).toString(CryptoJS.enc.Hex).substr(0, difficulty) === cmp){
	    return x;
	}
	x++;
    }
}
function find_xs(problems, difficulty){
	var res = [];
	for(var i = 0; i < problems.length; i++){
		var temp = {	Seed: problems[i].Seed,
				   		Solution: find_x(difficulty, problems[i].Seed)
				   };
		res.push(temp);
	}
	return res;
}
function log(msg){
	 $("#msg").append("<br/><b>"+msg+"</b><br/>");
}
function append_result(str){
	$("#result").append(str);
}

$(document).ready(function(){
	if (window["WebSocket"]) {
      conn = new WebSocket("ws://{{$}}/ws");
      conn.onclose = function(evt) {	
         log("Connection closed.");
      };
      conn.onmessage = function(evt) {
      	// alert("Got response: " + evt.data);
       	var response = JSON.parse(evt.data);
       	if(response["Opcode"] == 1){
       		// alert("Problems is:" + response.Problems);
       		$('#problems')[0].innerHTML = response.Difficulty.Problems;
       		$('#zeroes')[0].innerHTML = response.Difficulty.Zeroes;

	      	var solution = find_xs(response.Problems, response.Difficulty.Zeroes);
	    	var request = { "Problems": solution, "Query": response.Query, "Opcode": 1};
		   	conn.send(JSON.stringify(request));
	    } else {
	        var endTime = Number(new Date().getTime()); // returns the number of MS since the epoch	
	        $("#search").removeAttr("disabled");
			$("#result").html("<br/><b>Result from server</b><br/>");
			append_result(response["Query"]+"<br/>");
			append_result("<br/><b>Time for solving the puzzle: </b><bPr/>");
			append_result((endTime - startTime)+" ms <br/>");
			append_result("<br/><b>The solution was</b><br/>");
			append_result('"'+response["Result"]+'" <br/>');
	    }
     }

	       	// $("#result").append("<br/>" + solution + "<br/>");

// "Hash": CryptoJS.SHA256(solution + "" + response["Seed"]).toString(CryptoJS.enc.Hex),
        /*append_result("<br/><b>The hash value was</b><br/>");
				append_result(response["Hash"]+" <br/>");
				*/
    } else {
        appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"));
    }

    var startTime;
    /* 
     * This search function is protected by proof-of-work.
     */
   /* $("input").click(function(e){	
	    var request = {"Opcode": e.data, "Query": $('#search_field').val()};
	    $("#search").attr("disabled", "disabled");
	    startTime = Number(new Date().getTime());
	    conn.send(JSON.stringify(request));  // log("Sending " + JSON.stringify(request));
    });*/
	$("input#b1").click(function(e){
  var request = {"Opcode": 1, "Query": ""};
  conn.send(JSON.stringify(request));  // log("Sending " + JSON.stringify(request));
 });
	$("input#b2").click(function(e){
  var request = {"Opcode": 2, "Query": ""};
  conn.send(JSON.stringify(request));  // log("Sending " + JSON.stringify(request));
 });
	$("input#b3").click(function(e){
  var request = {"Opcode": 3, "Query": ""};
  conn.send(JSON.stringify(request));  // log("Sending " + JSON.stringify(request));
 });
});	


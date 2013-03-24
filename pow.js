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
	// $("#result").append("<br/><b>"+x+"<br/>");
	// $("#result").append("<br/><b>Calc is "+CryptoJS.SHA256(str).toString(CryptoJS.enc.Hex)+"<br/>");
	// $("#result").append("<br/><b>Cmp part "+CryptoJS.SHA256(str).toString(CryptoJS.enc.Hex).substr(0,difficulty)+"<br/>");
	if(CryptoJS.SHA256(str).toString(CryptoJS.enc.Hex).substr(0, difficulty) === cmp){
	    // alert("seed: " + seed +  "\nfound " + str +" after " +i);	    
	    return x;

	}
	x++;
    }
}
function find_xs(problems, difficulty){
	var res = [];
	for(var i = 0; i < problems.length; i++){
		var temp = {Seed: problems[i].Seed, Solution: find_x(difficulty, problems[i].Seed)};
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
	        	var solution = find_xs(response["Problems"], response["Difficulty"]);
	        	$("#result").append("<br/>" + solution + "<br/>");
		    	var request = { "Problems": solution, 
		    					"Query": response.Query,
		    					"Hash": CryptoJS.SHA256(solution + "" + response["Seed"]).toString(CryptoJS.enc.Hex),
		    					"Opcode": 1};
		    	conn.send(JSON.stringify(request));
	    	} else {
	            var endTime = Number(new Date().getTime()); // returns the number of MS since the epoch	
				$("#result").html("<br/><b>Result from server</b><br/>");
				append_result(response["Query"]+"<br/>");
				append_result("<br/><b>Time for solving the puzzle</b><br/>");
				append_result((endTime - startTime)+" ms <br/>");
				append_result("<br/><b>The solution was</b><br/>");
				append_result('"'+response["Result"]+'" <br/>');
				append_result("<br/><b>The hash value was</b><br/>");
				append_result(response["Hash"]+" <br/>");
	    	}
        }
    } else {
        appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"));
    }

    var startTime;
    /* 
     * This search function is protected by proof-of-work.
     */
    $("#search").click(function(){
		startTime = Number(new Date().getTime());
	    var request = {"Opcode": 0, "Query": $('#search_field').val()};
	    log("Sending " + JSON.stringify(request));
	    conn.send(JSON.stringify(request));  
    });
});


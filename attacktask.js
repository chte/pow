importScripts('sha256.js');

onmessage = function (e){
 //recieved data from UI

 var response = { id: e.data.id, 
 				  solution: find_solution(e.data.problems, e.data.difficulty)
 				};
postMessage(response)

 //sending message to UI
}


/****************************************** POW logics ************************************/

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
function find_solution(problems, difficulty){
	var res = [];
	for(var i = 0; i < problems.length; i++){
		var temp = {	Seed: problems[i].Seed,
				   		Solution: find_x(difficulty, problems[i].Seed)
				   };
		res.push(temp);
	}
	return res;
}

window.onload = function() {
	//定时五分钟
	var requestFrame =window.requestAnimationFrame || window.mozRequestAnimationFrame || window.webkitRequestAnimationFrame 
                                || window.msRequestAnimationFrame;
	//canvas init
	var canvas = document.getElementById("canvas");
	var ctx = canvas.getContext("2d");

	//canvas dimensions
	var W = window.innerWidth;
	var H = window.innerHeight;
	canvas.width = W;
	canvas.height = H;

	//雪花的密度
	var mp = 100; //max particles
	var particles = [];
	for (var i = 0; i < mp; i++) {
		particles.push({
			x: Math.random() * W, //x-coordinate
			y: Math.random() * H, //y-coordinate
			//雪花半径
			r: Math.random() * 4 + 1, //radius
			d: Math.random() * mp //density
		})
	}

	//Lets draw the flakes
	function draw() {
		ctx.clearRect(0, 0, W, H);
		ctx.fillStyle = "rgba(255, 255, 255, 0.8)";
		/* ctx.fillStyle = "#FF0000";*/
		ctx.beginPath();
		for (var i = 0; i < mp; i++) {
			var p = particles[i];
			ctx.moveTo(p.x, p.y);
			ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2, true);
		}
		ctx.fill();
		update();
	    requestFrame.call(window,draw);
	}

	//Function to move the snowflakes
	//angle will be an ongoing incremental flag. Sin and Cos functions will be applied to it to create vertical and horizontal movements of the flakes
	var angle = 0;

	function update() {
		angle += 0.01;
		for (var i = 0; i < mp; i++) {
			var p = particles[i];
			//Updating X and Y coordinates
			//We will add 1 to the cos function to prevent negative values which will lead flakes to move upwards
			//Every particle has its own density which can be used to make the downward movement different for each flake
			//Lets make it more random by adding in the radius
			//飞舞的速度
			p.y += Math.cos(angle + p.d) * 0.02 + p.r / 2 * 0.5;
			p.x += Math.sin(angle) * 0.2;

			//Sending flakes back from the top when it exits
			//Lets make it a bit more organic and let flakes enter from the left and right also.
			if (p.x > W || p.x < 0 || p.y > H) {
				if (i % 3 > 0) //66.67% of the flakes
				{
					particles[i] = {
						x: Math.random() * W,
						y: -10,
						r: p.r,
						d: p.d
					};
				} else {
					//If the flake is exitting from the right
					if (Math.sin(angle) > 0) {
						//Enter fromth
						particles[i] = {
							x: -5,
							y: Math.random() * H,
							r: p.r,
							d: p.d
						};
					} else {
						//Enter from the right
						particles[i] = {
							x: W + 5,
							y: Math.random() * H,
							r: p.r,
							d: p.d
						};
					}
				}
			}
		}
	}

	//animation loop
	//下落速度
	draw();
}
let add = fn(x) {
	fn(y) {
		y + x;
	}
}

let addTwo = add(2);
addTwo(7)

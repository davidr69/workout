parseUrlParams = () => {
	let parts = window.location.href.split('?');
	if(parts.length === 2) {
		let arr = { };
		for(let item of parts[1].split('&')) {
			let kvp = item.split('=');
			if(kvp.length === 2) {
				arr[kvp[0]] = kvp[1];
			}
		}
		return arr;
	} else {
		return null;
	}
}

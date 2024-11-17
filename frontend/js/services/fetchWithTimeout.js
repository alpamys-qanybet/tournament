export default function (url, options, timeout = 1000*60*60*24) {
	return Promise.race([
		fetch(url, options),
		new Promise((_, reject) =>
			setTimeout(() => reject({network: 'timeout'}), timeout)
		)
	]);
}
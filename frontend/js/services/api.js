import fetch from './fetchWithTimeout';

export const HTTP_STATUS_CODE_SUCCESS = 200;
export const HTTP_STATUS_CODE_CREATED = 201;
export const HTTP_STATUS_BAD_REQUEST = 400;
export const HTTP_STATUS_CODE_UNAUTHORIZED = 401;
export const HTTP_STATUS_CODE_UNPROCESSABLE_ENTITY = 422;

let site_url = window.location.host;

if (!(site_url.includes("http://") || site_url.includes("https://"))) {
	site_url = window.location.protocol + "//" + site_url;
}

export const SITE_URL = site_url;
const BASE_URL = `${SITE_URL}/api`;

const Api = {
	getTeamList: (callback, callbackErr) => {
		let status = HTTP_STATUS_CODE_SUCCESS;
		fetch(`${BASE_URL}/teams`)
		.then((response) => {
			status = response.status;
			return response;
		})
		.then((response) => response.json())
		.then((responseJson) => {
			return {
				status,
				data: responseJson,
			};
		})
		.then(callback)
		.catch((err) => {
			callbackErr(err);
		});
	},

	addTeam: (name, callback, callbackErr) => {
		let status = HTTP_STATUS_CODE_SUCCESS;
		fetch(`${BASE_URL}/teams`, {
			method: "POST",
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				name: name, 
			})
		})
		.then((response) => {
			status = response.status;
			return response;
		})
		.then((response) => response.json())
		.then((responseJson) => {
			return {
				status,
				data: responseJson,
			};
		})
		.then(callback)
		.catch((err) => {
			callbackErr(err);
		});
	},
	
	generateTeams: (callback, callbackErr) => {
		let status = HTTP_STATUS_CODE_SUCCESS;
		fetch(`${BASE_URL}/teams/generate`, {
			method: "POST",
		})
		.then((response) => {
			status = response.status;
			return response;
		})
		.then((response) => response.json())
		.then((responseJson) => {
			return {
				status,
				data: responseJson,
			};
		})
		.then(callback)
		.catch((err) => {
			callbackErr(err);
		});
	},

	getDivisions: (callback, callbackErr) => {
		let status = HTTP_STATUS_CODE_SUCCESS;
		fetch(`${BASE_URL}/divisions`)
		.then((response) => {
			status = response.status;
			return response;
		})
		.then((response) => response.json())
		.then((responseJson) => {
			return {
				status,
				data: responseJson,
			};
		})
		.then(callback)
		.catch((err) => {
			callbackErr(err);
		});
	},

	prepareDivision: (callback, callbackErr) => {
		let status = HTTP_STATUS_CODE_SUCCESS;
		fetch(`${BASE_URL}/divisions/prepare`, {
			method: "POST",
		})
		.then((response) => {
			status = response.status;
			return response;
		})
		.then((response) => response.json())
		.then((responseJson) => {
			return {
				status,
				data: responseJson,
			};
		})
		.then(callback)
		.catch((err) => {
			callbackErr(err);
		});
	},

	startDivision: (callback, callbackErr) => {
		let status = HTTP_STATUS_CODE_SUCCESS;
		fetch(`${BASE_URL}/divisions/start`, {
			method: "POST",
		})
		.then((response) => {
			status = response.status;
			return response;
		})
		.then((response) => response.json())
		.then((responseJson) => {
			return {
				status,
				data: responseJson,
			};
		})
		.then(callback)
		.catch((err) => {
			callbackErr(err);
		});
	},

	GetPlayoffs: (callback, callbackErr) => {
		let status = HTTP_STATUS_CODE_SUCCESS;
		fetch(`${BASE_URL}/playoff`)
		.then((response) => {
			status = response.status;
			return response;
		})
		.then((response) => response.json())
		.then((responseJson) => {
			return {
				status,
				data: responseJson,
			};
		})
		.then(callback)
		.catch((err) => {
			callbackErr(err);
		});
	},

	preparePlayoff: (stage, callback, callbackErr) => {
		let status = HTTP_STATUS_CODE_SUCCESS;
		fetch(`${BASE_URL}/playoff/prepare`, {
			method: "POST",
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				stage: stage, 
			})
		})
		.then((response) => {
			status = response.status;
			return response;
		})
		.then((response) => response.json())
		.then((responseJson) => {
			return {
				status,
				data: responseJson,
			};
		})
		.then(callback)
		.catch((err) => {
			callbackErr(err);
		});
	},

	startPlayoff: (stage, callback, callbackErr) => {
		let status = HTTP_STATUS_CODE_SUCCESS;
		fetch(`${BASE_URL}/playoff/start`, {
			method: "POST",
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				stage: stage, 
			})
		})
		.then((response) => {
			status = response.status;
			return response;
		})
		.then((response) => response.json())
		.then((responseJson) => {
			return {
				status,
				data: responseJson,
			};
		})
		.then(callback)
		.catch((err) => {
			callbackErr(err);
		});
	},

	cleanup: (callback, callbackErr) => {
		let status = HTTP_STATUS_CODE_SUCCESS;
		fetch(`${BASE_URL}/cleanup`, {
			method: "POST",
			// empty body, but cleanup is an operation, get is kind of for getting data in the meaning
		})
		.then((response) => {
			status = response.status;
			return response;
		})
		.then((response) => response.json())
		.then((responseJson) => {
			return {
				status,
				data: responseJson,
			};
		})
		.then(callback)
		.catch((err) => {
			callbackErr(err);
		});
	},
}

export default Api;
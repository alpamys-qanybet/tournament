import React, {Component} from 'react';
import LoadingOverlay from '../components/LoadingOverlay';
import Api, {
	HTTP_STATUS_CODE_SUCCESS,
	SITE_URL,
} from "../services/api";

class Cleanup extends Component {
	constructor(props) {
		super(props);
		this.state = {
			loading: false,
		};
	}

	_cleanup = () => {
		this.setState({
			loading: true,
		}, ()=> {
			Api.cleanup((result)=> {
				if (result.status === HTTP_STATUS_CODE_SUCCESS) {
					this.setState({
						loading: false,
					});

					window.location.replace(
						SITE_URL+"/teams",
					);
					return;
				}

				this.setState({
					loading: false,
				});
			}, (err)=> {
				this.setState({
					loading: false,
				});
			});
		});
	};

	render() {
		const {
			loading,
		} = this.state;

		return (
			<div>
				{loading && <LoadingOverlay/>}
				<div style={{
					paddingLeft: "20px",
					paddingRight: "20px",
				}}>
					<h5 style={{
						marginTop: "20px",
					}}>Cleanup</h5>
					<div style={{
						display: "flex",
						justifyContent: "flex-end",
						alignItems: "center",
					}}>
						<button className="btn btn-danger" onClick={this._cleanup}>Cleanup</button>
					</div>
				</div>
			</div>
		)
	}
}
export default Cleanup;
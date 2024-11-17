import React, {Component} from 'react';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';

import LoadingOverlay from '../components/LoadingOverlay';
import Api, {
	HTTP_STATUS_CODE_SUCCESS,
	HTTP_STATUS_CODE_CREATED,
	HTTP_STATUS_BAD_REQUEST,
	HTTP_STATUS_CODE_UNPROCESSABLE_ENTITY,
} from "../services/api";
import {defaultErrorMsg} from "../consts";


class Teams extends Component {
	constructor(props) {
		super(props);
		this.state = {
			loading: true,
			mode: "list", // list, add
			list: [],
			name: null, // adding team name
			isError: false,
			errMsg: "",
		};
	}

	componentDidMount() {
		this._fetch();
	}

	_fetch = () => {
		this.setState({
			loading: true,
		}, ()=> {
			Api.getTeamList((result)=> {

				if (result.status === HTTP_STATUS_CODE_SUCCESS) {
					this.setState({
						loading: false,
						list: result.data,
					});
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

	_add = () => {
		this.setState({
			isError: false,
			errMsg: "",
		});

		const {
			name,
		} = this.state;

		let nameIsNull = false;
		if (name) {
			if (name.trim().length === 0) {
				nameIsNull = true;
			}
		} else {
			nameIsNull = true;
		}

		if (nameIsNull) {
			this.setState({
				isError: true,
				errMsg: "Team name is required",
			});

			return;
		}

		this.setState({
			loading: true,
			isError: false,
			errMsg: "",
		}, () => {
			Api.addTeam(name, (result) => {
				if (result.status === HTTP_STATUS_CODE_CREATED) {
					this.setState({
						mode: "list",
						name: null,
					}, () => {
						this._fetch();
					});
					return;
				}

				if (result.status === HTTP_STATUS_BAD_REQUEST) {
					if (result.data.err === 'invalid_body_params') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Invalid body params",
						});
						return;
					}
				}
				
				if (result.status === HTTP_STATUS_CODE_UNPROCESSABLE_ENTITY) {
					
					if (result.data.err === 'create_team_failure_name_is_required') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Team name is required",
						});
						return;
					}

					if (result.data.err === 'max_16_teams_allowed') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Max 16 teams are allowed",
						});
						return;
					}
				}

				if (result.data.err) {
					this.setState({
						loading: false,
						isError: true,
						errMsg: result.data.err,
					});
					return;
				}

				this.setState({
					loading: false,
				});
			}, (err) => {
				this.setState({
					loading: false,
				});
			});
		});
	};

	_generate = () => {
		this.setState({
			loading: true,
		}, ()=> {
			Api.generateTeams((result)=> {
				if (result.status === HTTP_STATUS_CODE_SUCCESS) {
					this._fetch();
					return;
				}
				
				if (result.status === HTTP_STATUS_CODE_UNPROCESSABLE_ENTITY) {
					
					if (result.data.err === 'generation_only_allowed_into_empty_table') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Generation only allowed into empty table",
						});
						return;
					}
				}

				if (result.data.err) {
					this.setState({
						loading: false,
						isError: true,
						errMsg: result.data.err,
					});
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

	_renderList = () => {
		const {
			list,
			isError,
			errMsg,
		} = this.state;

		return (
			<div>
				<h5 style={{
					marginTop: "20px",
				}}>Teams</h5>
				<table className="table" style={{
					width: "100%",
					// backgroundColor: "red",
				}}>
					<thead>
						<tr>
							<th scope="col">#</th>
							<th scope="col">ID</th>
							<th scope="col">Name</th>
						</tr>
					</thead>
					<tbody>
						{list.map((t, i) => <tr key={"team-item_"+i+"_"+t.id} style={{
							verticalAlign: "baseline",
						}}>
							<th>{i + 1}</th>
							<td>{t.id}</td>
							<td>{t.name}</td>
						</tr>)}
					</tbody>
				</table>
				{list.length === 0 && <div>
					Empty, no teams
				</div>}
				{list.length === 16 && <div className="alert alert-primary" role="alert" style={{
					marginTop: "30px",
				}}>
					If you didn't start a division and more playoffs, then it is time, you can cleanup all data and start from scratch if you want
				</div>}
				{isError && (
					<div className="alert alert-danger" role="alert" style={{
						marginTop: "30px",
					}}>{errMsg || defaultErrorMsg}</div>
				)}
				{list.length < 16 && <div style={{
					marginTop: "30px",
					display: "flex",
					justifyContent: "flex-end",
				}}>
					{list.length === 0 && <button className="btn btn-success" onClick={this._generate}>Generate UCL Football teams</button>}
					<button className="btn btn-primary" onClick={()=> {
						this.setState({
							mode: "add",
							name: null,
							isError: false,
							errMsg: "",
						});
					}} style={{
						marginLeft: list.length === 0 ? "12px" : "0",
					}}>Add</button>
				</div>}
			</div>
		);
	}

	_renderAdd = () => {
		const {
			name,
			isError,
			errMsg,
		} = this.state;

		return (
			<div>
				<h5 style={{
					marginTop: "20px",
				}}>Add a team</h5>
				<div className="mt-3"> 
					<InputGroup className="mb-3">
						<InputGroup.Text id="add-team-name">Team name</InputGroup.Text>
						<Form.Control
							placeholder="Team name"
							aria-describedby="add-team-name"
							value={name} onChange={(e) => {
								e.preventDefault();
								this.setState({
									name: e.target.value,
								});
							}}
						/>
					</InputGroup>
				</div>

				{isError && (
					<div className="alert alert-danger" role="alert" style={{
						marginTop: "30px",
					}}>{errMsg || defaultErrorMsg}</div>
				)}

				<div style={{
					display: "flex",
					justifyContent: "flex-end",
					alignItems: "center",
				}}>
					<button className="btn btn-light" style={{
						borderWidth: "1px",
						borderColor: "#aeaeae",
					}}
					onClick={()=> {
						this.setState({
							mode: "list",
							name: null,
							isError: false,
							errMsg: "",
						});
					}}>Cancel</button>
					<button className="btn btn-primary" style={{
						marginLeft: "12px",
					}}
					onClick={this._add}>Add</button>
				</div>
			</div>
		);
	}

	render() {
		const {
			loading,
			mode,
		} = this.state;

		return (
			<div>
				{loading && <LoadingOverlay/>}
				<div style={{
					paddingLeft: "20px",
					paddingRight: "20px",
				}}>
					{mode === "list" && this._renderList()}
					{mode === "add" && this._renderAdd()}
				</div>
			</div>
		)
	}
}
export default Teams;
import React, {Component} from 'react';
import _, { map } from 'lodash';
import LoadingOverlay from '../components/LoadingOverlay';
import Api, {
	HTTP_STATUS_CODE_SUCCESS,
	HTTP_STATUS_BAD_REQUEST,
	HTTP_STATUS_CODE_UNPROCESSABLE_ENTITY,
} from "../services/api";
import {defaultErrorMsg} from "../consts";

class Playoff extends Component {

	teamMap = new Map();

	constructor(props) {
		super(props);
		this.state = {
			loading: true,
			quarter: {
				name:    "Quarter-Final",
				teams: [],
				matches: [],
			},
			semi: {
				name:    "Semi-Final",
				teams:   [],
				matches: [],
			},
			final: {
				name:    "Final",
				teams:   [],
				matches: [],
			},
			mode: "", // "prepare-quarter", // start-quarter, prepare-semi, start-semi, prepare-final, start-final, ended
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
			isError: false,
			errMsg: "",
		}, ()=> {
			Api.GetPlayoffs((result)=> {

				if (result.status === HTTP_STATUS_CODE_SUCCESS) {
					for (let v of result.data.quarter.teams) {
						this.teamMap.set(v.id, v.name);
					}

					this.setState({
						loading: false,
						quarter: result.data.quarter,
						semi: result.data.semi,
						final: result.data.final,
						mode: result.data.stage,
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

	_prepare = (stage) => {
		this.setState({
			loading: true,
			isError: false,
			errMsg: "",
		}, ()=> {
			Api.preparePlayoff(stage, (result)=> {
				if (result.status === HTTP_STATUS_CODE_SUCCESS) {
					this._fetch();
					return;
				}

				// invalid_body_params
				// division_is_not_started
				// playoff_quarter_is_already_prepared
				// playoff_quarter_is_already_started
				// playoff_semi_is_already_prepared
				// playoff_semi_is_already_started
				// playoff_final_is_already_prepared
				// playoff_final_is_already_started

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


					if (result.data.err === 'division_is_not_started') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Division is not started",
						});
						return;
					}

					if (result.data.err === 'playoff_quarter_is_already_prepared') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Quarter-Final Stage is already prepared",
						});
						return;
					}

					if (result.data.err === 'playoff_quarter_is_already_started') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Quarter-Final Stage is already started",
						});
						return;
					}

					if (result.data.err === 'playoff_semi_is_already_prepared') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Semi-Final Stage is already prepared",
						});
						return;
					}

					if (result.data.err === 'playoff_semi_is_already_started') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Semi-Final Stage is already started",
						});
						return;
					}

					if (result.data.err === 'playoff_final_is_already_prepared') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Final Stage is already prepared",
						});
						return;
					}

					if (result.data.err === 'playoff_final_is_already_started') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Final Stage is already started",
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

 	_start = (stage) => {
		this.setState({
			loading: true,
			isError: false,
			errMsg: "",
		}, ()=> {
			Api.startPlayoff(stage, (result)=> {
				if (result.status === HTTP_STATUS_CODE_SUCCESS) {
					this._fetch();
					return;
				}

				// invalid_body_params
				// playoff_quarter_is_not_prepared
				// playoff_quarter_is_already_started
				// playoff_semi_is_not_prepared
				// playoff_semi_is_already_started
				// playoff_final_is_not_prepared
				// playoff_final_is_already_started


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

					if (result.data.err === 'playoff_quarter_is_not_prepared') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Quarter-Final Stage is not prepared",
						});
						return;
					}

					if (result.data.err === 'playoff_quarter_is_already_started') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Quarter-Final Stage is already started",
						});
						return;
					}

					if (result.data.err === 'playoff_semi_is_not_prepared') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Semi-Final Stage is not prepared",
						});
						return;
					}

					if (result.data.err === 'playoff_semi_is_already_started') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Semi-Final Stage is already started",
						});
						return;
					}

					if (result.data.err === 'playoff_final_is_not_prepared') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Final Stage is not prepared",
						});
						return;
					}

					if (result.data.err === 'playoff_final_is_already_started') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Play-off Final Stage is already started",
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

	_renderPlayoff = (stage) => {
		const {
			quarter,
			semi,
			final,
		} = this.state;

		let teams = [];
		let matches = [];
		let stageTitle = "";
		if (stage === "quarter") {
			stageTitle = quarter.name;
			teams = _.cloneDeep(quarter.teams);
			matches = _.cloneDeep(quarter.matches);
		} else if (stage === "semi") {
			stageTitle = semi.name;
			teams = _.cloneDeep(semi.teams);
			matches = _.cloneDeep(semi.matches);
		} else if (stage === "final") {
			stageTitle = final.name;
			teams = _.cloneDeep(final.teams);
			matches = _.cloneDeep(final.matches);
		}

		return (
			<div>
				<h5 style={{
					marginTop: "70px",
				}}>{stageTitle} Stage</h5>
				<table className="table" style={{
					width: "100%",
				}}>
					<thead>
						<tr>
							<th scope="col">#</th>
							<th scope="col">ID</th>
							<th scope="col">Name</th>
						</tr>
					</thead>
					<tbody>
						{teams.map((t, i) => <tr key={stage+"team-item_"+i+"_"+t.id} style={{
							verticalAlign: "baseline",
						}}>
							<th>{i + 1}</th>
							<td>{t.id}</td>
							<td>{t.name}</td>
						</tr>)}
					</tbody>
				</table>
				{teams.length === 0 && <div>
					Empty, no teams
				</div>}

				{matches.length > 0 && <div>
					<h5 style={{
						marginTop: "40px",
					}}>{stageTitle} Matches</h5>
					
					<table className="table" style={{
						width: "100%",
					}}>
						<thead>
							<tr>
								<th scope="col">#</th>
								<th scope="col">ID</th>
								<th scope="col">Title</th>
								<th scope="col">Score</th>
								<th scope="col">Winner</th>
							</tr>
						</thead>
						<tbody>
							{matches.map((m, i) => {
								let winnerName = "";

								if (this.teamMap.has(m.winner_id)) {
									winnerName = this.teamMap.get(m.winner_id);
								}

								return (
									<tr key={"match-item_"+i+"_"+m.id} style={{
										verticalAlign: "baseline",
									}}>
										<th>{i + 1}</th>
										<td>{m.id}</td>
										<td>{m.name}</td>
										<td>{_.isNull(m.first_team_score) ? "" : (m.first_team_score + " - " + m.second_team_score)}</td>
										<td>{winnerName}</td>
									</tr>
								);
							})}
						</tbody>
					</table>
				</div>}
			</div>
		);
	};

	render() {
		const {
			loading,
			mode,
			isError,
			errMsg,
		} = this.state;

		return (
			<div>
				{loading && <LoadingOverlay/>}
				<div style={{
					paddingLeft: "20px",
					paddingRight: "20px",
					paddingBottom: "8px",
				}}>
					<h5 style={{
						marginTop: "20px",
					}}>Play-off</h5>
					{this._renderPlayoff("quarter")}
					{this._renderPlayoff("semi")}
					{this._renderPlayoff("final")}

					{isError && (
						<div className="alert alert-danger" role="alert" style={{
							marginTop: "30px",
						}}>{errMsg || defaultErrorMsg}</div>
					)}

					{mode === "prepare-quarter" && <div style={{
						display: "flex",
						justifyContent: "flex-end",
						alignItems: "center",
					}}>
						<button className="btn btn-primary" onClick={()=> {
							this._prepare("quarter");
						}}>Prepare Quarter-Final Matches</button>
					</div>}

					{mode === "prepare-semi" && <div style={{
						display: "flex",
						justifyContent: "flex-end",
						alignItems: "center",
					}}>
						<button className="btn btn-primary" onClick={()=> {
							this._prepare("semi");
						}}>Prepare Semi-Final Matches</button>
					</div>}

					{mode === "prepare-final" && <div style={{
						display: "flex",
						justifyContent: "flex-end",
						alignItems: "center",
					}}>
						<button className="btn btn-primary" onClick={()=> {
							this._prepare("final");
						}}>Prepare Final Match</button>
					</div>}

					{mode === "start-quarter" && <div style={{
						display: "flex",
						justifyContent: "flex-end",
						alignItems: "center",
					}}>
						<button className="btn btn-primary" onClick={()=> {
							this._start("quarter");
						}}>Start Quarter-Final Matches</button>
					</div>}

					{mode === "start-semi" && <div style={{
						display: "flex",
						justifyContent: "flex-end",
						alignItems: "center",
					}}>
						<button className="btn btn-primary" onClick={()=> {
							this._start("semi");
						}}>Start Semi-Final Matches</button>
					</div>}

					{mode === "start-final" && <div style={{
						display: "flex",
						justifyContent: "flex-end",
						alignItems: "center",
					}}>
						<button className="btn btn-primary" onClick={()=> {
							this._start("final");
						}}>Start Final Match</button>
					</div>}

					{mode === "ended" && <div className="alert alert-primary" role="alert" style={{
						marginTop: "30px",
					}}>
						We have a winner, you can cleanup all data and start from scratch if you want
					</div>}
				</div>
			</div>
		)
	}
}
export default Playoff;
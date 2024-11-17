import React, {Component} from 'react';
import _ from 'lodash';
import LoadingOverlay from '../components/LoadingOverlay';
import Api, {
	HTTP_STATUS_CODE_SUCCESS,
	HTTP_STATUS_CODE_UNPROCESSABLE_ENTITY,
} from "../services/api";
import {defaultErrorMsg} from "../consts";

class Divisions extends Component {
	constructor(props) {
		super(props);
		this.state = {
			loading: true,
			divisionA: {
				name: "A",
				teams: [],
				matches: [],
			},
			divisionB: {
				name: "B",
				teams: [],
				matches: [],
			},
			mode: "prepare", // start ended
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
			Api.getDivisions((result)=> {

				if (result.status === HTTP_STATUS_CODE_SUCCESS) {
					let a = _.cloneDeep(result.data.division_a);
					for (let i in a.teams) {
						if (!a.teams[i].hasOwnProperty("wins")) {
							a.teams[i].wins = 0;
						}
						if (!a.teams[i].hasOwnProperty("draws")) {
							a.teams[i].draws = 0;
						}
						if (!a.teams[i].hasOwnProperty("loses")) {
							a.teams[i].loses = 0;
						}
						if (!a.teams[i].hasOwnProperty("goals_scored")) {
							a.teams[i].goals_scored = 0;
						}
						if (!a.teams[i].hasOwnProperty("goals_conceded")) {
							a.teams[i].goals_conceded = 0;
						}
						if (!a.teams[i].hasOwnProperty("goal_diff")) {
							a.teams[i].goal_diff = 0;
						}
						if (!a.teams[i].hasOwnProperty("points")) {
							a.teams[i].points = 0;
						}
					}

					let b = _.cloneDeep(result.data.division_b);
					for (let i in b.teams) {
						if (!b.teams[i].hasOwnProperty("wins")) {
							b.teams[i].wins = 0;
						}
						if (!b.teams[i].hasOwnProperty("draws")) {
							b.teams[i].draws = 0;
						}
						if (!b.teams[i].hasOwnProperty("loses")) {
							b.teams[i].loses = 0;
						}
						if (!b.teams[i].hasOwnProperty("goals_scored")) {
							b.teams[i].goals_scored = 0;
						}
						if (!b.teams[i].hasOwnProperty("goals_conceded")) {
							b.teams[i].goals_conceded = 0;
						}
						if (!b.teams[i].hasOwnProperty("goal_diff")) {
							b.teams[i].goal_diff = 0;
						}
						if (!b.teams[i].hasOwnProperty("points")) {
							b.teams[i].points = 0;
						}
					}

					let mode = "ended"; // division ended
					if (a.teams.length === 0) { // division not prepared, maybe there not enough teams or maybe not activated divisions
						mode = "prepare"; // generate divisions from teams and schedule matches <=== prepare division
					} else {
						let totalMatches = a.teams[0].wins + a.teams[0].draws + a.teams[0].loses;
						if (totalMatches === 0) {
							mode = "start"; // generate match scores <=== start a division
						}
					}

					this.setState({
						loading: false,
						mode: mode,
						divisionA: a,
						divisionB: b,
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

	_prepare = () => {
		this.setState({
			loading: true,
			isError: false,
			errMsg: "",
		}, ()=> {
			Api.prepareDivision((result)=> {
				if (result.status === HTTP_STATUS_CODE_SUCCESS) {
					this._fetch();
					return;
				}

				if (result.status === HTTP_STATUS_CODE_UNPROCESSABLE_ENTITY) {
					
					if (result.data.err === 'must_have_16_teams_to_prepare_divisions') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Must have 16 teams to prepare divisions",
						});
						return;
					}

					if (result.data.err === 'division_is_already_prepared') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Divisions are already prepared",
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

	_start = () => {
		this.setState({
			loading: true,
			isError: false,
			errMsg: "",
		}, ()=> {
			Api.startDivision((result)=> {
				if (result.status === HTTP_STATUS_CODE_SUCCESS) {
					this._fetch();
					return;
				}

				if (result.status === HTTP_STATUS_CODE_UNPROCESSABLE_ENTITY) {
					
					if (result.data.err === 'must_have_16_teams_to_start_divisions') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Must have 16 teams to start divisions",
						});
						return;
					}

					if (result.data.err === 'division_is_not_prepared') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Divisions are not prepared",
						});
						return;
					}

					if (result.data.err === 'division_is_already_started') {
						this.setState({
							loading: false,
							isError: true,
							errMsg: "Divisions are already started",
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

	_renderDivision = (division) => {
		const {
			divisionA,
			divisionB,
			mode,
		} = this.state;

		let list = [];
		let matches = [];
		if (division === "A") {
			list = _.cloneDeep(divisionA.teams);
			matches = _.cloneDeep(divisionA.matches);
		} else if (division === "B") {
			list = _.cloneDeep(divisionB.teams);
			matches = _.cloneDeep(divisionB.matches);
		}

		return (
			<div>
				<h5 style={{
					marginTop: "20px",
				}}>Division - {division} Ranking Table</h5>
				<table className="table" style={{
					width: "100%",
					// backgroundColor: "red",
				}}>
					<thead>
						<tr>
							<th scope="col">#</th>
							<th scope="col">ID</th>
							<th scope="col">Name</th>
							<th scope="col">Matches Played</th>
							<th scope="col">Wins</th>
							<th scope="col">Draws</th>
							<th scope="col">Loses</th>
							<th scope="col">Goals scored</th>
							<th scope="col">Goals conceded</th>
							<th scope="col">Goal Diff</th>
							<th scope="col">Points</th>
						</tr>
					</thead>
					<tbody>
						{list.map((t, i) => <tr key={"team-item_"+i+"_"+t.id} style={{
							verticalAlign: "baseline",
						}}
						className={(mode === "ended" && i < 4)?"table-primary": ""}>
							<th>{i + 1}</th>
							<td>{t.id}</td>
							<td>{t.name}</td>
							<td>{t.wins+t.draws+t.loses}</td>
							<td>{t.wins}</td>
							<td>{t.draws}</td>
							<td>{t.loses}</td>
							<td>{t.goals_scored}</td>
							<td>{t.goals_conceded}</td>
							<td>{t.goal_diff}</td>
							<td>{t.points}</td>
						</tr>)}
					</tbody>
				</table>
				{list.length === 0 && <div>
					Empty, no teams
				</div>}

				{list.length > 0 && <div>
					<h5 style={{
						marginTop: "40px",
					}}>Division - {division} Matches</h5>
					
					<table className="table" style={{
						width: "100%",
					}}>
						<thead>
							<tr>
								<th scope="col">#</th>
								<th scope="col">ID</th>
								<th scope="col">Title</th>
								<th scope="col">Score</th>
							</tr>
						</thead>
						<tbody>
							{matches.map((m, i) => <tr key={"match-item_"+i+"_"+m.id} style={{
								verticalAlign: "baseline",
							}}>
								<th>{i + 1}</th>
								<td>{m.id}</td>
								<td>{m.name}</td>
								<td>{_.isNull(m.first_team_score) ? "" : (m.first_team_score + " - " + m.second_team_score)}</td>
							</tr>)}
						</tbody>
					</table>
				</div>}
			</div>
		);
	}

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
					}}>Divisions</h5>
					<div style={{
						marginTop: "70px",
					}}>
						{this._renderDivision("A")}
					</div>
					<div style={{
						marginTop: "70px",
					}}>
						{this._renderDivision("B")}
					</div>

					{isError && (
						<div className="alert alert-danger" role="alert" style={{
							marginTop: "30px",
						}}>{errMsg || defaultErrorMsg}</div>
					)}

					{mode === "prepare" && <div style={{
						display: "flex",
						justifyContent: "flex-end",
						alignItems: "center",
					}}>
						<button className="btn btn-primary" onClick={this._prepare}>Prepare</button>
					</div>}

					{mode === "start" && <div style={{
						display: "flex",
						justifyContent: "flex-end",
						alignItems: "center",
					}}>
						<button className="btn btn-primary" onClick={this._start}>Start</button>
					</div>}

					{mode === "ended" && <div className="alert alert-primary" role="alert" style={{
						marginTop: "30px",
					}}>
						If you didn't start playoffs, then it is time, you can cleanup all data and start from scratch if you want
					</div>}
				</div>
			</div>
		)
	}
}
export default Divisions;
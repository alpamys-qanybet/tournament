import React, {Component} from 'react';
import { Link, Outlet } from 'react-router-dom'; 

import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';

class Layout extends Component {
	
	menu = [{
		link: "/teams",
		code: "teams",
		title: "Teams",
	}, {
		link: "/divisions",
		code: "divisions",
		title: "Divisions",
	}, {
		link: "/playoff",
		code: "playoff",
		title: "Playoff",
	}, {
		link: "/cleanup",
		code: "cleanup",
		title: "Cleanup",
	}];

	constructor(props) {
		super(props);
		
		let m = this._getActiveMenu(); 
		this.state = {
			active: m,
		};
	}
	
	_getActiveMenu = () => {
		const url = window.location.pathname;

		let index = this.menu.map(e => e.link).indexOf(url);
		if (index == -1) {
			index = 0;
		}

		return this.menu[index].code
	}

	render() {
		const {
			active,
		} = this.state;

		return (
			<main>
				<Navbar className="bg-body-tertiary" bg="dark" data-bs-theme="dark" sticky="top">
					<Container style={{
						marginLeft: "10px",
						marginRight: "10px",
						width: "100%",
					}}>
						<Navbar.Brand as={Link} to="/" onClick={()=> {
							this.setState({
								active: this.menu[0].code,
							});
						}}>Tournament</Navbar.Brand>
						<Navbar.Toggle />
						<Navbar.Collapse className="justify-content-start">
							<Nav>
								{this.menu.map((m, i) => {
									return (
										<Nav.Link key={"menu-item_"+i}
										className={"" + (active === m.code ? "active" : "")}
										as={Link}
										to={m.link}
										onClick={()=> {
											this.setState({
												active: m.code,
											});
										}}>{m.title}</Nav.Link>
									);
								})}
							</Nav>
						</Navbar.Collapse>
					</Container>
				</Navbar>
				<Outlet/>
			</main>
		)
	}
}

export default Layout;
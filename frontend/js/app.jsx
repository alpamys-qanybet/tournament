import React from 'react';
import ReactDOM from 'react-dom/client';
import {
	BrowserRouter as Router,
	Route,
	Routes,
} from "react-router-dom";

import ScrollToTop from './components/ScrollToTop';

import Layout from './pages/Layout';
import Teams from "./pages/Teams";
import Divisions from "./pages/Divisions";
import Playoff from './pages/Playoff';
import Cleanup from './pages/Cleanup';

class App extends React.Component {

	constructor(props) {
		super(props);

		this.state = {
		};
	}

	render() {
		return (
			<Router>
				<ScrollToTop/>
				
				<Routes>
					<Route path={"/"} element={<Layout/>}>
						<Route index element={<Teams/>}/>
						<Route path="/teams" element={<Teams/>}/>
						<Route path="/divisions" element={<Divisions/>}/>
						<Route path="/playoff" element={<Playoff/>}/>
						<Route path="/cleanup" element={<Cleanup/>}/>
					</Route>
				</Routes>
			</Router>
		);
	}
}

const renderApp = () => {
	root.render(
		// <React.StrictMode> // calls ComponentDidMount twice
			<App/>
		// </React.StrictMode>
	);
};

const root = ReactDOM.createRoot(document.getElementById("root"));
renderApp();

import React, {Component} from 'react';

export default class LoadingOverlay extends Component {

	render() {
		return (
			<div style={{
				display: "flex",
				position: 'absolute',
				top: 0,
				left: 0,
				right: 0,
				bottom: 0,
				height: '100vh',
				width: '100vw',
				backgroundColor: 'rgba(219,219,219,0.5)',
				zIndex: 3000,
				justifyContent: 'center',
				alignItems: 'center',
			}}>
				<div style={{
					marginTop: "-100px", // "-50px"
					width: '100px',
					height: '100px',
					padding: 19,
					borderRadius: 13,
					alignContent: "center",
					justifyContent: 'center',
					backgroundColor: 'rgba(39,48,64,0.9)',
					color: "white",
				}}>
					Loading...
				</div>
			</div>
		);
	}
}

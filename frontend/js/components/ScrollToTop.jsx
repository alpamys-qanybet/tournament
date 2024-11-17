import { useEffect } from "react";
import { useLocation } from "react-router";

const ScrollToTop = ({ children }) => {
	const { pathname } = useLocation();

	useEffect(() => {
		window.scrollTo({ top:0, left:0, behavior: "instant"});
		var scrollableContainer = document.getElementById('scrollable-container');
		if (scrollableContainer) {
			scrollableContainer.scrollTop = 0;
		}
	}, [pathname]);

	return null;
};

export default ScrollToTop;
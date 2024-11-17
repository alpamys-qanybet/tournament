const esbuild = require("esbuild");

esbuild.context({
	entryPoints: ["frontend/js/app.jsx", "frontend/css/style.css"],
	outdir: "assets/app",
	bundle: true,
	minify: true, // prod
	loader: {
		'.png': 'file',
		'.jpg': 'file',
	},
})
.then((ctx) => ctx.watch())
.then(() => console.log("⚡ Build complete! ⚡"))
.then(() => console.log("⚡ Watching... ⚡"))
.catch(() => process.exit(1));
const path = require('path');

const HTMLWebpackPlugin = require('html-webpack-plugin');
const MiniCSSExtractPlugin = require('mini-css-extract-plugin');
const CleanWebpackPlugin = require('clean-webpack-plugin');

const webpack = require('webpack');
const autoprefixer = require('autoprefixer');

module.exports = (env, argv) => {

	const isProduction = argv.mode === 'production';
	console.log("Building for production: ", isProduction);

	return {
		entry: "./src/index.jsx",
		output: {
			path: path.resolve(__dirname, 'dist'),
			publicPath: "/jstn/"
		},
		devtool: isProduction ? false : 'inline-source-map',
		devServer: {
			publicPath: "/jstn/",
			hot: true
		},
		module: {
			rules: [
				{
					test: /\.(js|jsx)$/,
					exclude: /node_modules/,
					use: [
						{ loader: "babel-loader" },
						{ loader: "eslint-loader" }
					]
				},
				{
					test: /\.html$/,
					use: [
						{ loader: "html-loader", options: { minimize: isProduction === true } }
					]
				},
				{
					test: /\.less$/,
					use: [
						{
							loader: isProduction ? MiniCSSExtractPlugin.loader : "style-loader"
						},
						{
							loader: "css-loader",
							options: {
								sourceMap: !isProduction,
								modules: true,
								localIdentName: "[local]__[hash:base64:5]"
							}
						},
						{
							loader: "postcss-loader",
							options: {
								ident: "postcss",
								plugins: [
									autoprefixer({
										browsers: [
											">1%",
											"last 2 major versions",
											"not ie < 9"
										],
										flexbox: 'no-2009'
									})
								]
							}
						},
						{ loader: "less-loader" }
					]
				}
			]
		},
		plugins: [
			new CleanWebpackPlugin(['dist']),
			new HTMLWebpackPlugin({
				template: "./src/index.html",
				filename: "./index.html"
			}),
			new MiniCSSExtractPlugin({
				filename: "[name].css",
				chunkFilename: "[id].css"
			}),
			new webpack.NamedModulesPlugin(),
			new webpack.HotModuleReplacementPlugin()
		]
	}
};
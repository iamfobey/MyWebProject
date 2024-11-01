const path = require("path");
const HTMLWebpackPlugin = require("html-webpack-plugin");

module.exports = {
    entry: "./src/index.js",
    output: {
        filename: "bundle.js",
        path: path.resolve("dist"),
        publicPath: "/",
    },
    module: {
        rules:[
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: "babel-loader"
            },
            {
                test: /\.html$/,
                use: "html-loader"
            },
        ],
    },
    devServer: {
        port: 3000,
        open: true,
        historyApiFallback: true,
        proxy: [
            {
                context: ["api"],
                target: "http://localhost:8080",
                pathRewrite: { '^/api': '' },
            }
        ]
    },
    plugins: [
        new HTMLWebpackPlugin({
            template: "./public/index.html"
        }),
    ]
}
const path = require("path");
const HTMLWebpackPlugin = require("html-webpack-plugin");

module.exports = {
    entry: "./www/src/index.js",
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'www/dist'),
        publicPath: '/'
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
        compress: true,
        host: '0.0.0.0',
        historyApiFallback: true,
    },
    plugins: [
        new HTMLWebpackPlugin({
            template: "./www/public/index.html"
        }),
    ]
}
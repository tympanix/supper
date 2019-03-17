const { DefinePlugin } = require("webpack");
const path = require("path");
const merge = require("webpack-merge");
const FaviconsWebpackPlugin = require("favicons-webpack-plugin");
const UglifyJSPlugin = require("uglifyjs-webpack-plugin");
const common = require("./webpack.common.js");

module.exports = merge(common, {
  mode: "production",
  devtool: "source-map",
  plugins: [
    new UglifyJSPlugin({
      sourceMap: true
    }),
    new DefinePlugin({
      "process.env.NODE_ENV": JSON.stringify("production")
    }),
    new FaviconsWebpackPlugin({
      logo: path.join(__dirname, "../favicon.png"),
      prefix: "favicon/",
      title: "Supper"
    })
  ]
});

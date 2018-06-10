const path = require('path');
const webpack = require('webpack');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

var devMode = process.env.NODE_ENV === "development"

module.exports = {
  entry: './web/js/index.js',
  output: {
    path: path.join(__dirname, 'web', 'static'),
    publicPath: "/static/",
    filename: 'bundle.js'
  },
  plugins: [
    new MiniCssExtractPlugin({
      filename: "styles.css",
    }),
  ],
  module: {
    rules: [
      {
        test: /\.scss$/,
        use: [
          devMode ? 'style-loader' : MiniCssExtractPlugin.loader,
          'css-loader',
          'sass-loader',
        ],
      }, {
        test: /\.js$/,
        loader: 'babel-loader',
        query: {
          presets: ['env', 'react']
        }
      }, {
        test: /\.svg$/,
        loader: 'file-loader',
        options: {
          publicPath: './img',
          outputPath: './img',
        }
      }
    ]
  },
  stats: {
    colors: true
  },
  devServer: {
    contentBase: 'web',
    inline: true,
    host: "0.0.0.0",
    proxy: {
      '/api': {
        target: 'http://localhost:5670',
        secure: false
      }
    },
    historyApiFallback: {
      index: 'index.html'
    }
  },
  devtool: 'source-map'
};

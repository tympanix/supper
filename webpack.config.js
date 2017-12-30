const path = require('path');
const webpack = require('webpack');
const ExtractTextPlugin = require("extract-text-webpack-plugin");

const extractSass = new ExtractTextPlugin({
  filename: "styles.css",
  disable: process.env.NODE_ENV === "development"
});

module.exports = {
  entry: './web/js/index.js',
  output: {
    path: path.join(__dirname, 'web', 'static'),
    publicPath: "/static/",
    filename: 'bundle.js'
  },
  module: {
    rules: [{
      test: /\.scss$/,
      use: extractSass.extract({
        use: [{
          loader: "css-loader"
        }, {
          loader: "sass-loader"
        }],
        // use style-loader in development
        fallback: "style-loader"
      })
    }, {
      test: /\.js$/,
      loader: 'babel-loader',
      query: {
        presets: ['es2015', 'react']
      }
    }]
  },
  plugins: [
    extractSass
  ],
  stats: {
    colors: true
  },
  devServer: {
    contentBase: 'web',
    inline: true,
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

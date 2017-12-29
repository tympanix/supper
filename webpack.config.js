var path = require('path');
var webpack = require('webpack');

 module.exports = {
     entry: './web/js/index.js',
     output: {
         path: path.join(__dirname, 'web', 'static'),
         publicPath: "/static/",
         filename: 'bundle.js'
     },
     module: {
         loaders: [
             {
                 test: /\.js$/,
                 loader: 'babel-loader',
                 query: {
                     presets: ['es2015', 'react']
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
       proxy: {
         '/api': {
           target: 'http://localhost:5670',
           secure: false
         }
       }
    },
    devtool: 'source-map'
 };

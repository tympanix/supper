const merge = require('webpack-merge');
const common = require('./webpack.common.js');

module.exports = merge(common, {
  mode: 'development',
  devtool: 'inline-source-map',
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
});

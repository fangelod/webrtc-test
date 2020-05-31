const path = require('path');
const webpack = require('webpack');

module.exports = options => ({
  entry: options.entry,
  output: Object.assign(
    {
      path: path.resolve(process.cwd(), 'dist'),
      filename: 'testrtc.js',
      library: ['TestRTC'],
      libraryTarget: 'umd',
    },
    options.output
  ),
  module: {
    rules: options.module.rules.concat([{
      test: /\.(eot|svg|ttf|woff|woff2)$/,
      use: [{
        loader: 'file-loader',
        options: {
          name: './fonts/[name].[ext]'
        }
      }]  
    }, {
      test: /\.mp3$/,
      use: [{
        loader: 'file-loader',
        options: {
          name: './audio/[name].[ext]'
        }
      }]
    }, {
      test: /\.css$/,
      use: ['style-loader', {
        loader: 'css-loader',
        options: {
          url: false
        }
      }]
    }, {
      test: /\.html$/,
      use: ['html-loader'],
    }, {
      test: /\.json$/,
      use: ['json-loader'],
      exclude: /node_modules/,
    }, {
      test: /\.(mp4|webm)$/,
      use: [{
        loader: 'url-loader',
        options: {
          limit: 10000,
        }
      }]
    }])
  },
  plugins: options.plugins.concat([
    new webpack.ProvidePlugin({
      fetch: 'exports-loader?self.fetch!whatwg-fetch'
    }),
    new webpack.ContextReplacementPlugin(
      /\.\/locale$/,
      'empty-module',
      false,
      /js$/
    ),
    new webpack.NamedModulesPlugin()
  ]),
  resolve: {
    modules: [
      path.resolve(process.cwd(), 'app'),
      path.resolve(process.cwd(), 'node_modules'),
      path.resolve(process.cwd(), 'app/containers')
    ],
    extensions: ['.js', '.jsx', '.react.js'],
    mainFields: ['browser', 'jsnext:main', 'main']
  },
  devtool: options.devtool,
  target: 'web',
  performance: options.performance || {}
});
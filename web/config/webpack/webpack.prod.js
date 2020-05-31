const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const DashboardPlugin = require('webpack-dashboard/plugin');
const UglifyJsPlugin = require('uglifyjs-webpack-plugin');

const pluginList = [
  new webpack.DefinePlugin({
    'process.env': {
      NODE_ENV: '\'production\'',
    }
  }),
  new UglifyJsPlugin({
    sourceMap: true
  }),
  new HtmlWebpackPlugin({
    template: 'app/index.html',
    minify: {
      removeComments: true,
      collapseWhitespace: true,
      removeRedundantAttributes: true,
      useShortDoctype: true,
      removeEmptyAttributes: true,
      removeStyleLinkTypeAttributes: true,
      keepClosingSlash: true,
      minifyJS: true,
      minifyCSS: true,
      minifyURLs: true
    },
    inject: 'head',
    favicon: 'resources/webrtc-logo.png'
  })
];

module.exports = env => require('./webpack.base')({
  entry: ['idempotent-babel-polyfill', path.resolve(process.cwd(), 'app/main.js')],
  output: {
    path: path.resolve(process.cwd(), 'dist'),
    filename: 'chat.js',
    library: ['NovaChat'],
    libraryTarget: 'umd',
  },
  module: {
    rules: [{
      test: /\.(js|jsx)$/,
      include: [
        path.resolve(process.cwd(), 'app')
      ],
      use: [{
        loader: 'babel-loader',
        options: {
          presets: [
            '@babel/preset-env',
            '@babel/preset-react'
          ],
          plugins: [
            'transform-react-remove-prop-types',
            '@babel/plugin-transform-react-constant-elements',
            '@babel/plugin-transform-react-inline-elements',
            'dynamic-import-node',
            '@babel/plugin-proposal-class-properties'
          ]
        }
      }]
    }, {
      test: /\.(jpg|png|gif)$/,
      use: [{
        loader: 'file-loader',
        options: {
          name: './images/[name].[ext]',
        }
      }]
    }]
  },
  plugins: env && env.dashboard ? pluginList.concat(new DashboardPlugin()) : pluginList,
  performance: {
    assetFilter: assetFilename =>
      !/(\.map$)|(^(main\.|favicon\.))/.test(assetFilename),
  }
});
const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const DashboardPlugin = require('webpack-dashboard/plugin');

module.exports = () => require('./webpack.base')({
  entry: ['idempotent-babel-polyfill', path.resolve(process.cwd(), 'app/main.js')],
  output: {
    path: path.resolve(process.cwd(), 'dist'),
    filename: 'testrtc.js',
    library: ['TestRTC'],
    libraryTarget: 'umd'
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
            '@babel/plugin-transform-react-constant-elements',
            '@babel/plugin-transform-react-inline-elements',
            'dynamic-import-node',
            '@babel/plugin-proposal-class-properties'
          ]
        }
      }],
    }, {
      test: /\.(jpg|png|gif)$/,
      use: [{
        loader: 'file-loader',
        options: {
          name: './images/[name].[ext]'
        }
      }, {
        loader: 'image-webpack-loader',
        options: {
          gifsicle: {
            progressive: true,
            interlaced: false,
            optimizationLevel: 7
          },
          mozjpeg: {
            progressive: true,
            interlaced: false
          },
          optipng: {
            progressive: true,
            interlaced: false,
            optimizationLevel: 7
          },
          pngquant: {
            quality: '65-90',
            speed: 4
          }
        }
      }]
    }]
  },
  plugins: [
    new DashboardPlugin(),
    new HtmlWebpackPlugin({
      template: 'app/index.html',
      inject: 'head',
      favicon: 'resources/webrtc-logo.png'
    })
  ],
  devtool: 'eval-source-map',
});
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const webpack = require('webpack');
const path = require('path');
const url = require('url');
const UglifyJSPlugin = require('uglifyjs-webpack-plugin');
const OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin');

require('dotenv').config();

const isDev = process.env.NODE_ENV === 'development';

function getScriptHost() {
  const hostStr = (process.env.WEBPACK_DEV_SERVER_HOST !== undefined) ?
    process.env.WEBPACK_DEV_SERVER_HOST : 'http://localhost:8081';
  const parsed = url.parse(hostStr);
  const publicPath = `${parsed.href}_goserve/assets`;
  return {
    publicPath,
    url: url.parse(hostStr),
  };
}

const plugins = isDev ? [
  new webpack.HotModuleReplacementPlugin(),
] : [
  new ExtractTextPlugin({
    filename: 'css/[name].css',
  }),
  new OptimizeCssAssetsPlugin({
    assetNameRegExp: /\.css$/,
    cssProcessorOptions: {
      discardComments: {
        removeAll: true,
      }
    },
  }),
  new UglifyJSPlugin(),
];

const sassRule = isDev ? {
  test: /\.scss$/,
  use: [
    // creates style nodes from JS strings
    'style-loader',
    // translates CSS into CommonJS
    'css-loader',
    // compiles Sass to CSS
    'sass-loader',
  ],
} : {
  test: /\.scss$/,
  use: ExtractTextPlugin.extract({
    fallback: 'style-loader',
    use: [
      'css-loader',
      'sass-loader',
    ],
  }),
};

const externals = {
};

const scriptHost = getScriptHost();

module.exports = {
  entry: {
    'app': [
      'babel-polyfill',
      './assets/src/js/app.js',
    ],
  },
  output: {
    path: path.resolve(__dirname, 'assets/dist'),
    publicPath: !isDev ? '' : scriptHost.publicPath,
    filename: 'js/[name].js',
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        include: path.join(__dirname, 'dist', 'src', 'js'),
        use: [
          'react-hot-loader',
          'babel-loader',
        ],
      },
      sassRule,
    ],
  },
  externals,
  plugins,
  devServer: {
    hot: true, // this enables hot reload
    // hotOnly: true, // do not reload browser if hot reload failed
    inline: true, // use inline method for hmr
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, PATCH, OPTIONS',
      'Access-Control-Allow-Headers': 'X-Requested-With, content-type, Authorization',
    },
    host: scriptHost.url.hostname,
    port: scriptHost.url.port,
    contentBase: path.join(__dirname, 'public'),
    watchOptions: {
      poll: false,
    },
  },
};

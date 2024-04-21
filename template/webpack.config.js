const HtmlWebpackPlugin = require("html-webpack-plugin");
const path = require("path");

module.exports = {
  entry: {
    login: "./src/login.tsx", // メインページ用エントリポイント
    error: "./src/error.tsx", // エラーページ用エントリポイント
  },
  output: {
    path: path.resolve(__dirname, "dist"),
    filename: "[name].bundle.js",
  },
  plugins: [
    new HtmlWebpackPlugin({
      title: "Go-IdP Login",
      filename: "login.html", // メインページの HTML ファイル名
      template: "src/login.html", // 使用する HTML テンプレート
      chunks: ["login"], // メインページ用のチャンクのみを含める
    }),
    new HtmlWebpackPlugin({
      title: "Go-IdP Error",
      filename: "error.html", // エラーページの HTML ファイル名
      template: "src/error.html", // 使用する HTML テンプレート
      chunks: ["error"], // エラーページ用のチャンクのみを含める
    }),
  ],
  resolve: {
    extensions: [".tsx", ".ts", ".js"],
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: "babel-loader",
        exclude: /node_modules/,
      },
    ],
  },
  mode: "development",
};

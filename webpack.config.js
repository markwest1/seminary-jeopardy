module.exports = {
  entry: "./public/app/main.ts",
  mode: "development",
  output: { filename: "./bundle.js" },
  resolve: {
    extensions: [".webpack.js", ".web.js", ".ts", ".tsx", ".js"],
  },
  module: {
    rules: [
      { test: /\.tsx?$/, loader: "ts-loader" },
      { test: /\.jade$/, loader: "jade" },
    ],
  },
};

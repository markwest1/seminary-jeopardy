module.exports = {
  entry: "./public/app/main.ts",
  output: { filename: "./public/dist/bundle.js" },
  resolve: {
    extensions: [".webpack.js", ".web.js", ".ts", ".tsx", ".js"],
  },
  module: {
    rules: [{ test: /\.tsx?$/, loader: "ts-loader" }],
  },
};

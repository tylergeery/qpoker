const path = require('path');
const glob = require('glob');

module.exports = {
    mode: process.env.WEBPACK_MODE || 'development',

    // Enable sourcemaps for debugging webpack's output.
    devtool: "source-map",

    resolve: {
        // Add '.ts' and '.tsx' as resolvable extensions.
        extensions: [".ts", ".tsx", ".js"]
    },

    module: {
        rules: [
            {
                test: /\.ts(x?)$/,
                exclude: /node_modules/,
                use: [
                    {
                        loader: "ts-loader"
                    }
                ]
            },
            // All output '.js' files will have any sourcemaps re-processed by 'source-map-loader'.
            {
                enforce: "pre",
                test: /\.js$/,
                loader: "source-map-loader"
            }
        ]
    },

    entry: Object.assign(
        // static files
        {
            main: './src/main.tsx',
            common: './src/common.tsx',
        },

        // dynamic game files
        glob.sync('./src/games/*.tsx')
            .reduce(
                (entries, gameFile) => {
                    entries['games/' + path.parse(gameFile).name] = gameFile;
                    return entries;
                },
                {}
            )
    ),
    output: {
        filename: '[name].bundle.js',
        path: path.resolve(__dirname, 'assets/js'),
    },

    // When importing a module whose path matches one of the following, just
    // assume a corresponding global variable exists and use that instead.
    // This is important because it allows us to avoid bundling all of our
    // dependencies, which allows browsers to cache those libraries between builds.
    // externals: {
    //     "react": "React",
    //     "react-dom": "ReactDOM"
    // }
};
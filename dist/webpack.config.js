const path = require('path');
const glob = require('glob');
const webpack = require('webpack');
const Dotenv = require('dotenv-webpack');

// 直接加载 .env 文件到 process.env
require('dotenv').config({ path: path.resolve(__dirname, '.env') });

module.exports = {
    entry: glob.sync('./*.ts').filter(file => !file.endsWith('.d.ts')).reduce((entries, file) => {
        const name = path.basename(file, '.ts');
        entries[name] = `./${file}`;
        return entries;
    }, {}),
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, './')
    },
    mode: 'production',
    module: {
        rules: [{
            test: /\.tsx?$/,
            use: 'ts-loader',
            exclude: /node_modules/
        }]
    },
    resolve: {
        extensions: ['.ts']
    },
    plugins: [
        new Dotenv({
            path: path.resolve(__dirname, '.env'),  // 显式指定 .env 路径
            systemvars: true,  // 允许系统环境变量覆盖 .env 文件
            silent: false,     // 如果 .env 不存在，显示警告
            defaults: false    // 不使用 .env.defaults 文件
        }),
        new webpack.DefinePlugin({
            '__BSZ_API_DOMAIN__': JSON.stringify(
                process.env.BSZ_API_DOMAIN || 'https://busuanzi.xxdevops.cn'
            )
        })
    ]
}
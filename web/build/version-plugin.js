const fs = require('fs');
const path = require('path');
const childProcess = require('child_process');


class VersionInfoPlugin {
  apply(compiler) {
    compiler.hooks.done.tap('VersionInfoPlugin', (stats) => {
      try {
        const pkg = require('../package.json');//定义版本,每个版本更新需要修改,中间版本已后端为主
        
        // 安全获取 git 信息：使用 execFileSync 替代 execSync，防止命令注入
        const gitPath = process.env.GIT_PATH || 'git';
        let branch = '';
        let commitId = '';
        let commitTime = '';
        
        try {
          branch = childProcess.execFileSync(
            gitPath,
            ['rev-parse', '--abbrev-ref', 'HEAD'],
            { encoding: 'utf8', timeout: 5000, maxBuffer: 1024 * 1024 }
          ).trim();
        } catch (e) {
          // git 命令失败时静默处理，不影响构建流程
        }
        
        try {
          commitId = childProcess.execFileSync(
            gitPath,
            ['rev-parse', 'HEAD'],
            { encoding: 'utf8', timeout: 5000, maxBuffer: 1024 * 1024 }
          ).trim();
        } catch (e) {
          // git 命令失败时静默处理，不影响构建流程
        }
        
        try {
          commitTime = childProcess.execFileSync(
            gitPath,
            ['log', '-1', '--format=%cd'],
            { encoding: 'utf8', timeout: 5000, maxBuffer: 1024 * 1024 }
          ).trim();
        } catch (e) {
          // git 命令失败时静默处理，不影响构建流程
        }
        
        const notes = '版本基于后端版本定义';
        const versionInfo = {
          version: pkg.version,
          branch,
          commitId,
          commitTime,
          notes,
          buildTime: new Date().toLocaleString()
        };
        
        // 获取webpack配置的输出目录，默认为dist
        const outputPath = stats.compilation.outputOptions.path || path.join(process.cwd(), 'dist');
        
        // 验证输出路径安全性：确保路径在项目目录内，防止路径遍历攻击
        const projectRoot = process.cwd();
        const resolvedPath = path.resolve(outputPath);
        if (!resolvedPath.startsWith(projectRoot)) {
          throw new Error('输出路径不在项目目录内');
        }
        
        // 确保目录存在
        if (!fs.existsSync(outputPath)) {
          fs.mkdirSync(outputPath, { recursive: true });
        }
        
        fs.writeFileSync(
          path.join(outputPath, 'version.json'),
          JSON.stringify(versionInfo, null, 2)
        );
        
        console.log('版本信息文件已生成:', path.join(outputPath, 'version.json'));
      } catch (e) {
        // 生产环境不输出详细错误信息，避免泄露系统路径
        const isProduction = process.env.NODE_ENV === 'production';
        if (isProduction) {
          console.error('生成版本信息失败');
        } else {
          console.error('生成版本信息失败:', e.message);
        }
      }
    });
  }
}

module.exports = VersionInfoPlugin;
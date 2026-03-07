/**
 * 格式化路由标题（兼容空值，避免match报错）
 * @param {string|undefined|null} title - 原始标题（可能包含${变量}）
 * @param {object|undefined|null} now - 路由对象（包含params/query）
 * @returns {string} 格式化后的标题
 */
export const fmtTitle = (title, now) => {
  // 1. 标题为空时直接返回空字符串，避免后续match报错
  if (!title || typeof title !== 'string') {
    return '';
  }

  const reg = /\$\{(.+?)\}/;
  const reg_g = /\$\{(.+?)\}/g;
  const result = title.match(reg_g);

  // 2. 没有匹配到${变量}，直接返回原标题
  if (!result || result.length === 0) {
    return title;
  }

  // 3. 确保now是对象，避免访问params/query时报错
  const routeData = now || { params: {}, query: {} };
  const { params = {}, query = {} } = routeData;

  // 4. 遍历替换变量，增加空值校验
  result.forEach((item) => {
    // 4.1 校验item匹配结果，避免[1]访问报错
    const matchResult = item.match(reg);
    if (!matchResult || !matchResult[1]) {
      return; // 匹配失败则跳过当前项
    }

    const key = matchResult[1];
    // 4.2 优先取params，再取query，都没有则保留原变量名
    const value = params[key] || query[key] || item;
    // 4.3 替换标题中的变量
    title = title.replace(item, value);
  });

  return title;
};
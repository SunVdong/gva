// src/hooks/useDict.js
import { ref, onMounted, unref } from 'vue';
import { getDict, showDictLabel } from '@/utils/dictionary'; // 导入项目现有工具函数

/**
 * 字典Hooks（适配项目现有dictionary utils）
 * @param type 字典类型（必填）
 * @param options 字典查询参数（depth/value）
 * @param autoLoad 是否自动加载（默认true）
 * @returns 响应式字典数据、加载状态、标签展示方法等
 */
export function useDict(
  type,
  options = { depth: 0, value: null },
  autoLoad = true
) {
  // 响应式字典数据
  const dictData = ref([]);
  // 加载状态
  const dictLoading = ref(false);
  // 错误信息
  const dictError = ref(null);

  /**
   * 加载字典数据
   * @param updateOptions 可选，更新查询参数
   */
  const loadDict = async (updateOptions) => {
    if (!type) {
      dictError.value = '字典类型不能为空';
      return [];
    }

    const finalOptions = { ...options, ...(updateOptions || {}) };
    dictLoading.value = true;
    dictError.value = null;

    try {
      // 调用项目现有getDict方法
      const data = await getDict(type, finalOptions);
      dictData.value = data;
      return data;
    } catch (error) {
      dictError.value = `加载字典失败: ${error.message}`;
      console.error('useDict load error:', error);
      return [];
    } finally {
      dictLoading.value = false;
    }
  };

  /**
   * 获取字典标签（简化调用）
   * @param code 字典值
   * @param keyCode 自定义值字段（默认'value'）
   * @param valueCode 自定义标签字段（默认'label'）
   */
  const getDictLabel = (
    code,
    keyCode = 'value',
    valueCode = 'label'
  ) => {
    return showDictLabel(unref(dictData), code, keyCode, valueCode);
  };

  /**
   * 刷新字典数据
   * @param updateOptions 可选，更新查询参数
   */
  const refreshDict = (updateOptions) => {
    return loadDict(updateOptions);
  };

  // 自动加载字典
  onMounted(() => {
    if (autoLoad) {
      loadDict();
    }
  });

  return {
    dictData,      // 响应式字典数据
    dictLoading,   // 加载状态
    dictError,     // 错误信息
    loadDict,      // 手动加载字典
    getDictLabel,  // 获取字典标签
    refreshDict    // 刷新字典
  };
}

/**
 * 批量加载字典的Hooks
 * @param dictConfig 字典配置数组 [{ type: 'xxx', options: { depth: 0, value: null } }]
 * @param autoLoad 是否自动加载（默认true）
 */
export function useMultiDict(
  dictConfig,
  autoLoad = true
) {
  // 存储所有字典数据 { [type]: [] }
  const dictMap = ref({});
  // 批量加载状态
  const dictLoading = ref(false);
  // 错误信息 { [type]: string }
  const dictErrors = ref({});

  /**
   * 加载所有字典
   */
  const loadAllDict = async () => {
    if (!dictConfig || dictConfig.length === 0) return;

    dictLoading.value = true;
    dictErrors.value = {};

    try {
      // 批量加载字典
      const promises = dictConfig.map(async (item) => {
        try {
          const data = await getDict(item.type, item.options || { depth: 0, value: null });
          dictMap.value[item.type] = data;
          dictErrors.value[item.type] = null;
        } catch (error) {
          dictErrors.value[item.type] = `加载${item.type}失败: ${error.message}`;
          console.error(`useMultiDict load ${item.type} error:`, error);
          dictMap.value[item.type] = [];
        }
      });

      await Promise.all(promises);
    } finally {
      dictLoading.value = false;
    }
  };

  /**
   * 获取指定字典的标签
   * @param type 字典类型
   * @param code 字典值
   * @param keyCode 自定义值字段
   * @param valueCode 自定义标签字段
   */
  const getDictLabelByType = (
    type,
    code,
    keyCode = 'value',
    valueCode = 'label'
  ) => {
    const dict = unref(dictMap)[type] || [];
    return showDictLabel(dict, code, keyCode, valueCode);
  };

  /**
   * 刷新指定字典
   * @param type 字典类型
   * @param options 可选，更新查询参数
   */
  const refreshDictByType = async (
    type,
    options
  ) => {
    const config = dictConfig.find(item => item.type === type);
    if (!config) {
      console.warn(`useMultiDict: 未找到字典配置 ${type}`);
      return;
    }

    dictLoading.value = true;
    try {
      const data = await getDict(type, options || config.options);
      dictMap.value[type] = data;
      dictErrors.value[type] = null;
    } catch (error) {
      dictErrors.value[type] = `刷新${type}失败: ${error.message}`;
      console.error(`useMultiDict refresh ${type} error:`, error);
    } finally {
      dictLoading.value = false;
    }
  };

  // 自动加载所有字典
  onMounted(() => {
    if (autoLoad) {
      loadAllDict();
    }
  });

  return {
    dictMap,              // 所有字典数据 { [type]: [] }
    dictLoading,          // 批量加载状态
    dictErrors,           // 错误信息
    loadAllDict,          // 加载所有字典
    getDictLabelByType,   // 根据类型获取标签
    refreshDictByType     // 刷新指定字典
  };
}
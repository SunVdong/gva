<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="景点名称" prop="name">
          <el-input v-model="searchInfo.name" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="searchInfo.status" placeholder="请选择">
            <el-option
              v-for="item in dictMap.scenic_status"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="创建日期" prop="createdAt">
          <el-date-picker
            v-model="searchInfo.createdAt"
            type="daterange"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            unlink-panels
          >
          </el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit"
            >查询</el-button
          >
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog"
          >新增</el-button
        >

        <el-button
          icon="delete"
          style="margin-left: 10px"
          :disabled="!multipleSelection.length"
          @click="onDelete"
          >删除</el-button
        >
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column
          align="left"
          label="景点名称"
          prop="name"
          width="120"
        />
        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{
            formatDate(scope.row.CreatedAt)
          }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="状态"
          min-width="120"
          show-overflow-tooltip
        >
          <template #default="scope">{{
            getDictLabelByType("scenic_status", scope.row.status)
          }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="280">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateSysExportTemplateFunc(scope.row)"
              >编辑</el-button
            >
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-drawer
      v-model="dialogFormVisible"
      size="60%"
      :before-close="closeDialog"
      :title="type === 'create' ? '添加' : '修改'"
      :show-close="false"
      destroy-on-close
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === "create" ? "添加" : "修改" }}</span>
          <div>
            <el-button @click="closeDialog">取 消</el-button>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="elFormRef"
        :model="formData"
        label-position="right"
        :rules="rule"
        label-width="100px"
        v-loading="aiLoading"
        element-loading-text="小淼正在思考..."
      >
        <el-form-item label="业务库" prop="dbName">
          <template #label>
            <el-tooltip
              content="注：需要提前到db-list自行配置多数据库，如未配置需配置后重启服务方可使用。若无法选择，请到config.yaml中设置disabled:false，选择导入导出的目标库。"
              placement="bottom"
              effect="light"
            >
              <div>
                业务库 <el-icon><QuestionFilled /></el-icon>
              </div>
            </el-tooltip>
          </template>
          <el-select
            v-model="formData.dbName"
            clearable
            @change="dbNameChange"
            placeholder="选择业务库"
          >
            <el-option
              v-for="item in dbList"
              :key="item.aliasName"
              :value="item.aliasName"
              :label="item.aliasName"
              :disabled="item.disable"
            >
              <div>
                <span>{{ item.aliasName }}</span>
                <span style="float: right; color: #8492a6; font-size: 13px">{{
                  item.dbName
                }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="需用到的表" prop="tables">
          <el-select
            multiple
            v-model="tables"
            clearable
            placeholder="使用AI的情况下请选择"
          >
            <el-option
              v-for="item in tableOptions"
              :key="item.tableName"
              :label="item.tableName"
              :value="item.tableName"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="AI帮写:" prop="ai">
          <div class="relative w-full">
            <el-input
              type="textarea"
              v-model="prompt"
              :clearable="true"
              :rows="5"
              placeholder="试试描述你要做的导出功能让AI帮你完成，在此之前请选择你需要导出的表所在的业务库，如不做选择，则默认使用gva库"
            />
            <el-button
              class="absolute bottom-2 right-2"
              type="primary"
              @click="autoExport"
              ><el-icon><ai-gva /></el-icon>帮写</el-button
            >
          </div>
        </el-form-item>

        <el-form-item label="表名称:" clearable prop="tableName">
          <div class="w-full flex gap-4">
            <el-select
              v-model="formData.tableName"
              class="flex-1"
              filterable
              placeholder="请选择表"
            >
              <el-option
                v-for="item in tableOptions"
                :key="item.tableName"
                :label="item.tableName"
                :value="item.tableName"
              />
            </el-select>
            <el-button
              :disabled="!formData.tableName"
              type="primary"
              @click="getColumnFunc(true)"
              ><el-icon><ai-gva /></el-icon>自动补全</el-button
            >
            <el-button
              :disabled="!formData.tableName"
              type="primary"
              @click="getColumnFunc(false)"
              >自动生成模板</el-button
            >
          </div>
        </el-form-item>

        <el-form-item label="模板名称:" prop="name">
          <el-input
            v-model="formData.name"
            :clearable="true"
            placeholder="请输入模板名称"
          />
        </el-form-item>

        <el-form-item label="模板标识:" prop="templateID">
          <el-input
            v-model="formData.templateID"
            :clearable="true"
            placeholder="模板标识为前端组件需要挂在的标识属性"
          />
        </el-form-item>

        <el-tabs v-model="activeName">
          <el-tab-pane label="自动构建" name="auto" class="pt-2">
            <el-form-item label="关联条件:">
              <div
                v-for="(join, key) in formData.joinTemplate"
                :key="key"
                class="flex gap-4 w-full mb-2"
              >
                <el-select v-model="join.joins" placeholder="请选择关联方式">
                  <el-option label="LEFT JOIN" value="LEFT JOIN" />
                  <el-option label="INNER JOIN" value="INNER JOIN" />
                  <el-option label="RIGHT JOIN" value="RIGHT JOIN" />
                </el-select>
                <el-input v-model="join.table" placeholder="请输入关联表" />
                <el-input
                  v-model="join.on"
                  placeholder="关联条件 table1.a = table2.b"
                />
                <el-button
                  type="danger"
                  icon="delete"
                  @click="() => formData.joinTemplate.splice(key, 1)"
                  >删除</el-button
                >
              </div>
              <div class="flex justify-end w-full">
                <el-button type="primary" icon="plus" @click="addJoin"
                  >添加条件</el-button
                >
              </div>
            </el-form-item>

            <el-form-item label="默认导出条数:">
              <el-input-number
                v-model="formData.limit"
                :step="1"
                :step-strictly="true"
                :precision="0"
              />
            </el-form-item>
            <el-form-item label="默认排序条件:">
              <el-input v-model="formData.order" placeholder="例:id desc" />
            </el-form-item>
            <el-form-item label="导出条件:">
              <div
                v-for="(condition, key) in formData.conditions"
                :key="key"
                class="flex gap-4 w-full mb-2"
              >
                <el-input
                  v-model="condition.from"
                  placeholder="需要从查询条件取的json key"
                />
                <el-input
                  v-model="condition.column"
                  placeholder="表对应的column"
                />
                <el-select
                  v-model="condition.operator"
                  placeholder="请选择查询条件"
                >
                  <el-option
                    v-for="item in typeSearchOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
                <el-button
                  type="danger"
                  icon="delete"
                  @click="() => formData.conditions.splice(key, 1)"
                  >删除</el-button
                >
              </div>
              <div class="flex justify-end w-full">
                <el-button type="primary" icon="plus" @click="addCondition"
                  >添加条件</el-button
                >
              </div>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="自定义SQL" name="sql" class="pt-2">
            <el-form-item label="导出SQL:" prop="sql">
              <el-input
                v-model="formData.sql"
                type="textarea"
                :rows="10"
                placeholder="请输入导出SQL语句，支持GORM命名参数模式，例如：SELECT * FROM sys_apis WHERE id = @id"
              />
            </el-form-item>
            <el-form-item label="导入SQL:" prop="importSql">
              <el-input
                v-model="formData.importSql"
                type="textarea"
                :rows="10"
                placeholder="请输入导入SQL语句，支持GORM命名参数模式，例如：INSERT INTO sys_apis (path, description ,api_group, method) VALUES (@path, @description, @api_group, @method)。参数名对应模板信息中的key。"
              />
            </el-form-item>
            <el-form-item label="导出条件:">
              此时导出条件的key必然为 condition =
              {key1:"value1",key2:"value2"}，这里需要和你传入sql语句@key占位符的key一致。
            </el-form-item>
          </el-tab-pane>
        </el-tabs>

        <el-form-item label="模板信息:" prop="templateInfo">
          <el-input
            v-model="formData.templateInfo"
            type="textarea"
            :rows="12"
            :clearable="true"
            :placeholder="templatePlaceholder"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createSysExportTemplate,
  deleteSysExportTemplate,
  deleteSysExportTemplateByIds,
  updateSysExportTemplate,
  findSysExportTemplate,
  getSysExportTemplateList,
} from "@/api/exportTemplate.js";
import { previewSQL } from "@/api/exportTemplate.js";

// 全量引入格式化工具 请按需保留
import { formatDate } from "@/utils/format";
import { useMultiDict } from "@/hooks/useDict";
import { ElMessage, ElMessageBox } from "element-plus";
import { ref, reactive } from "vue";
import { getDB, getTable, getColumn, llmAuto } from "@/api/autoCode";
import { VAceEditor } from "vue3-ace-editor";

import "ace-builds/src-noconflict/mode-vue";
import "ace-builds/src-noconflict/theme-github_dark";
import "ace-builds/src-noconflict/mode-sql";

defineOptions({
  name: "Scenic",
});

const templatePlaceholder = `模板信息格式：key标识数据库column列名称（在join模式下需要写为 table.column），value标识导出excel列名称，如key为数据库关键字或函数，请按照关键字的处理模式处理，当前以mysql为例，如下：
{
  "table_column1":"第一列",
  "table_column3":"第三列",
  "table_column4":"第四列",
  "\`rows\`":"我属于数据库关键字或函数",
}
如果使用是sql模式，您自行构建的sql的key就是需要写在json的key，例如您写了xxx as k1，那么模板信息中就写{"k1":"对应列名称"}
如果增加了JOINS导出key应该列为 {table_name1.table_column1:"第一列",table_name2.table_column2:"第二列"}
如果有重复的列名导出格式应为 {table_name1.table_column1 as key:"第一列",table_name2.table_column2 as key2:"第二列"}
JOINS模式下不支持导入
`;

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  name: "",
  tableName: "",
  dbName: "",
  templateID: "",
  templateInfo: "",
  limit: 0,
  order: "",
  conditions: [],
  joinTemplate: [],
  sql: "",
  importSql: "",
});

const activeName = ref("auto");

const prompt = ref("");
const tables = ref([]);

const typeSearchOptions = ref([
  {
    label: "=",
    value: "=",
  },
  {
    label: "<>",
    value: "<>",
  },
  {
    label: ">",
    value: ">",
  },
  {
    label: "<",
    value: "<",
  },
  {
    label: "LIKE",
    value: "LIKE",
  },
  {
    label: "BETWEEN",
    value: "BETWEEN",
  },
  {
    label: "NOT BETWEEN",
    value: "NOT BETWEEN",
  },
  {
    label: "IN",
    value: "IN",
  },
  {
    label: "NOT IN",
    value: "NOT IN",
  },
]);

const addCondition = () => {
  formData.value.conditions.push({
    from: "",
    column: "",
    operator: "",
  });
};

const addJoin = () => {
  formData.value.joinTemplate.push({
    joins: "LEFT JOIN",
    table: "",
    on: "",
  });
};

// 验证规则
const rule = reactive({
  name: [
    {
      required: true,
      message: "",
      trigger: ["input", "blur"],
    },
    {
      whitespace: true,
      message: "不能只输入空格",
      trigger: ["input", "blur"],
    },
  ],
  tableName: [
    {
      required: true,
      message: "",
      trigger: ["input", "blur"],
    },
    {
      whitespace: true,
      message: "不能只输入空格",
      trigger: ["input", "blur"],
    },
  ],
  templateID: [
    {
      required: true,
      message: "",
      trigger: ["input", "blur"],
    },
    {
      whitespace: true,
      message: "不能只输入空格",
      trigger: ["input", "blur"],
    },
  ],
  templateInfo: [
    {
      required: true,
      message: "",
      trigger: ["input", "blur"],
    },
    {
      whitespace: true,
      message: "不能只输入空格",
      trigger: ["input", "blur"],
    },
  ],
});
const elFormRef = ref();
const elSearchFormRef = ref();
const { dictMap, dictLoading, getDictLabelByType, refreshDictByType } =
  useMultiDict([{ type: "scenic_status", options: { depth: 0 } }]);
// =========== 表格控制部分 ===========
const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const tableData = ref([]);
const searchInfo = ref({
  createdAt: [],
  name: null,
  status: null,
});

const dbList = ref([]);
const tableOptions = ref([]);
const aiLoading = ref(false);

const getTablesCloumn = async () => {
  const tablesMap = {};
  const promises = tables.value.map(async (item) => {
    const res = await getColumn({
      businessDB: formData.value.dbName,
      tableName: item,
    });
    if (res.code === 0) {
      tablesMap[item] = res.data.columns;
    }
  });
  await Promise.all(promises);
  return tablesMap;
};

const autoExport = async () => {
  if (tables.value.length === 0) {
    ElMessage({
      type: "error",
      message: "请先选择需要参与导出的表",
    });
    return;
  }
  aiLoading.value = true;
  const tableMap = await getTablesCloumn();
  const aiRes = await llmAuto({
    prompt: prompt.value,
    tableMap: JSON.stringify(tableMap),
    mode: "autoExportTemplate",
  });
  aiLoading.value = false;
  if (aiRes.code === 0) {
    const aiData = JSON.parse(aiRes.data);
    formData.value.name = aiData.name;
    formData.value.tableName = aiData.tableName;
    formData.value.templateID = aiData.templateID;
    formData.value.templateInfo = JSON.stringify(aiData.templateInfo, null, 2);
    formData.value.joinTemplate = aiData.joinTemplate;
  }
};

const getDbFunc = async () => {
  const res = await getDB();
  if (res.code === 0) {
    dbList.value = res.data.dbList;
  }
};

getDbFunc();

const dbNameChange = () => {
  formData.value.tableName = "";
  formData.value.templateInfo = "";
  tables.value = [];
  getTableFunc();
};

const getTableFunc = async () => {
  const res = await getTable({ businessDB: formData.value.dbName });
  if (res.code === 0) {
    tableOptions.value = res.data.tables;
  }
  formData.value.tableName = "";
};
getTableFunc();
const getColumnFunc = async (aiFLag) => {
  if (!formData.value.tableName) {
    ElMessage({
      type: "error",
      message: "请先选择业务库及选择表后再进行操作",
    });
    return;
  }
  formData.value.templateInfo = "";
  aiLoading.value = true;
  const res = await getColumn({
    businessDB: formData.value.dbName,
    tableName: formData.value.tableName,
  });
  if (res.code === 0) {
    if (aiFLag) {
      const aiRes = await llmAuto({
        data: JSON.stringify(res.data.columns),
        mode: "exportCompletion",
      });
      if (aiRes.code === 0) {
        const aiData = JSON.parse(aiRes.data);
        aiLoading.value = false;
        formData.value.templateInfo = JSON.stringify(
          aiData.templateInfo,
          null,
          2
        );
        formData.value.name = aiData.name;
        formData.value.templateID = aiData.templateID;
        return;
      }
      ElMessage.warning("AI自动补全失败，已调整为逻辑填写");
    }

    // 把返回值的data.columns做尊换，制作一组JSON数据，columnName做key，columnComment做value
    const templateInfo = {};
    res.data.columns.forEach((item) => {
      templateInfo[item.columnName] = item.columnComment || item.columnName;
    });
    formData.value.templateInfo = JSON.stringify(templateInfo, null, 2);
  }
  aiLoading.value = false;
};

// 重置
const onReset = () => {
  searchInfo.value = {};
  getTableData();
};

// 搜索
const onSubmit = () => {
  page.value = 1;
  getTableData();
};

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val;
  getTableData();
};

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val;
  getTableData();
};

// 查询
const getTableData = async () => {
  const createdAt = searchInfo.value.createdAt || [];
  const table = await getSysExportTemplateList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value,
    createStartDate: createdAt[0] ?? null,
    createEndDate: createdAt[1] ?? null,
  });
  if (table.code === 0) {
    tableData.value = table.data.list;
    total.value = table.data.total;
    page.value = table.data.page;
    pageSize.value = table.data.pageSize;
  }
};

getTableData();

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {};

// 获取需要的字典 可能为空 按需保留
setOptions();

// 多选数据
const multipleSelection = ref([]);
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val;
};

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm("确定要删除吗?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    deleteSysExportTemplateFunc(row);
  });
};

// 多选删除
const onDelete = async () => {
  ElMessageBox.confirm("确定要删除吗?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(async () => {
    const ids = [];
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: "warning",
        message: "请选择要删除的数据",
      });
      return;
    }
    multipleSelection.value &&
      multipleSelection.value.map((item) => {
        ids.push(item.ID);
      });
    const res = await deleteSysExportTemplateByIds({ ids });
    if (res.code === 0) {
      ElMessage({
        type: "success",
        message: "删除成功",
      });
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--;
      }
      getTableData();
    }
  });
};

// 行为控制标记（弹窗内部需要增还是改）
const type = ref("");

// 更新行
const updateSysExportTemplateFunc = async (row) => {
  const res = await findSysExportTemplate({ ID: row.ID });
  type.value = "update";
  if (res.code === 0) {
    formData.value = res.data.resysExportTemplate;
    if (!formData.value.conditions) {
      formData.value.conditions = [];
    }
    if (!formData.value.joinTemplate) {
      formData.value.joinTemplate = [];
    }
    if (!formData.value.sql) {
      formData.value.sql = "";
    }
    if (!formData.value.importSql) {
      formData.value.importSql = "";
    }
    if (formData.value.sql || formData.value.importSql) {
      activeName.value = "sql";
    } else {
      activeName.value = "auto";
    }
    dialogFormVisible.value = true;
  }
};

// 删除行
const deleteSysExportTemplateFunc = async (row) => {
  const res = await deleteSysExportTemplate({ ID: row.ID });
  if (res.code === 0) {
    ElMessage({
      type: "success",
      message: "删除成功",
    });
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--;
    }
    getTableData();
  }
};

// 弹窗控制标记
const dialogFormVisible = ref(false);

// 打开弹窗
const openDialog = () => {
  type.value = "create";
  dialogFormVisible.value = true;
};

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false;
  formData.value = {
    name: "",
    tableName: "",
    templateID: "",
    templateInfo: "",
    limit: 0,
    order: "",
    conditions: [],
    joinTemplate: [],
    sql: "",
    importSql: "",
  };
  activeName.value = "auto";
};
// 弹窗确定
const enterDialog = async () => {
  // 判断 formData.templateInfo 是否为标准json格式 如果不是标准json 则辅助调整
  try {
    JSON.parse(formData.value.templateInfo);
  } catch (_) {
    ElMessage({
      type: "error",
      message: "模板信息格式不正确，请检查",
    });
    return;
  }

  const reqData = JSON.parse(JSON.stringify(formData.value));
  if (activeName.value === "sql") {
    reqData.conditions = [];
    reqData.joinTemplate = [];
    reqData.limit = 0;
    reqData.order = "";
  } else {
    reqData.sql = "";
    reqData.importSql = "";
  }

  for (let i = 0; i < reqData.conditions.length; i++) {
    if (
      !reqData.conditions[i].from ||
      !reqData.conditions[i].column ||
      !reqData.conditions[i].operator
    ) {
      ElMessage({
        type: "error",
        message: "请填写完整的导出条件",
      });
      return;
    }
    reqData.conditions[i].templateID = reqData.templateID;
  }

  for (let i = 0; i < reqData.joinTemplate.length; i++) {
    if (!reqData.joinTemplate[i].joins || !reqData.joinTemplate[i].on) {
      ElMessage({
        type: "error",
        message: "请填写完整的关联",
      });
      return;
    }
    reqData.joinTemplate[i].templateID = reqData.templateID;
  }

  elFormRef.value?.validate(async (valid) => {
    if (!valid) return;
    let res;
    switch (type.value) {
      case "create":
        res = await createSysExportTemplate(reqData);
        break;
      case "update":
        res = await updateSysExportTemplate(reqData);
        break;
      default:
        res = await createSysExportTemplate(reqData);
        break;
    }
    if (res.code === 0) {
      ElMessage({
        type: "success",
        message: "创建/更改成功",
      });
      closeDialog();
      getTableData();
    }
  });
};
</script>

<style></style>

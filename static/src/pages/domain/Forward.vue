<template>
  <div class="app">
    <el-table border :data="list" class="table">
      <el-table-column align="center" prop="type" label="类型" width="60">
        <template #default="item">
          <span v-if="item.row.type === 1">ipv4</span>
          <span v-else>ipv6</span>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="pattern" label="表达式"/>
      <el-table-column align="center" prop="priority" label="优先级"/>
      <el-table-column align="center" prop="ttl" label="ttl"/>
      <el-table-column align="center" prop="values" label="值">
        <template #default="item">
          <span :title="JSON.stringify(item.row.values)"
                v-text="$tool.splitText(JSON.stringify(item.row.values),25)"></span>
        </template>
      </el-table-column>
      <el-table-column align="center" fixed="right" label="操作" width="200">
        <template #default>
          <el-button type="primary" size="small">详情</el-button>
          <el-button type="warning" size="small">编辑</el-button>
          <el-button type="danger" size="small">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
        style="float: right"
        background
        layout="total, prev, pager, next "
        :hide-on-single-page="false"
        :total="total"
        @current-change="pageChange"
        :default-current-page="page"
        :default-page-size="size"
    />
  </div>
</template>
<script>
export default {
  name: "ForwardPage",
  data: () => {
    return {
      list: [],
      page: 1,
      size: 15,
      total: 0
    }
  },
  methods: {
    pageChange(page){
      this.page = page
      this.dnsSearch()
    },
    dnsSearch() {
      this.$api.dnsSearch({
        type: "1,28",
        page: this.page,
        size: this.size
      }).then(res => {
        let data = res.data
        if (!data.state) {
          this.$message.error(data.msg)
          return
        }
        this.total = data.attr.total
        this.list = data.data
      })
    }
  },
  mounted() {
    this.dnsSearch()
  }
}
</script>

<style scoped>
.app {
  width: 100%;
  border: 1px solid rgba(255, 255, 255, 0);
  box-sizing: border-box;
  height: 100%;
}

.table {
  margin-top: 100px;
  height: calc(100% - 160px);;
}
</style>
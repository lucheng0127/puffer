<template>
  <div class="login-cls">
    <div class="login-top">
      <el-icon class="icon-cls"><Ship /></el-icon>
      <span class="system-title">
        {{ systemName }}
      </span>
    </div>
    <div class="login-home">
      <div class="home-left">
        <div class="home-notice">
          <p class="flex-align-center">
            <el-icon class="icon-cls notice-title"><Connection /></el-icon>
            <span class="notice-title">温馨提示</span>
          </p>
          <el-divider />
          <p
            class="notice-info"
            v-for="(item, index) in noticeMsg"
            :key="index"
          >
            {{ item }}
          </p>
        </div>
      </div>
      <div class="home-right">
        <div class="login-box">
          <template v-if="loginFlag">
            <p class="login-title flex-align-center">
              账号已登录
              <el-icon style="color: green; margin-left: 8px">
                <CircleCheckFilled />
              </el-icon>
            </p>
            <p>{{ loginMsg }}</p>
            <p>
              <el-button
                type="primary"
                @click="handleLogout"
                class="login-button"
              >
                退出连接
              </el-button>
            </p>
          </template>
          <template v-else>
            <p class="login-title">账号登录</p>
            <div class="box-cls">
              <el-form
                label-width="auto"
                ref="loginFormRef"
                :model="loginForm"
                :rules="loginRules"
              >
                <el-form-item label="服务器地址" prop="target">
                  <el-input v-model="loginForm.target" />
                </el-form-item>
                <el-form-item label="用户名" prop="username">
                  <el-input v-model="loginForm.username" />
                </el-form-item>
                <el-form-item label="密码" prop="password">
                  <el-input
                    v-model="loginForm.password"
                    type="password"
                    show-password
                  />
                </el-form-item>
                <el-form-item label="密钥" prop="key">
                  <el-input
                    v-model="loginForm.key"
                    type="password"
                    show-password
                  />
                </el-form-item>
                <el-form-item label="日志级别" prop="logLevel">
                  <el-select v-model="loginForm.logLevel" placeholder="">
                    <el-option
                      v-for="item in logLeveOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
                <el-form-item>
                  <el-button
                    type="primary"
                    @click="handleSumit(loginFormRef)"
                    class="login-button"
                    :loading="loginLoading"
                    >连接</el-button
                  >
                </el-form-item>
              </el-form>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { reactive, Ref, ref } from "vue";
import { Ship, Connection, CircleCheckFilled } from "@element-plus/icons-vue";
import { FormInstance, FormRules } from "element-plus";
import { Launch, Logout, CheckAuthed } from "../../../wailsjs/go/main/App";
import { ElMessage } from "element-plus";

const systemName = ref("Virtuallan客户端");
const noticeMsg: Ref<string[]> = ref([
  "欢迎使用Virtuallan客户端，让您的系统接入更安全，体验更好！",
]);
const logLeveOptions = [
  { value: "info", label: "INFO" },
  { value: "debug", label: "DEBUG" },
];
const loginFlag = ref(false);
const loginMsg = ref("-");
const loginFormRef = ref();
const loginForm = reactive({
  target: "",
  username: "",
  password: "",
  key: "",
  logLevel: "",
});
const loginRules = reactive<FormRules>({
  target: [{ required: true, message: "请输入服务器地址", trigger: "blur" }],
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
  key: [{ required: true, message: "请输入密钥", trigger: "blur" }],
});

const loginLoading = ref(false);
const wsUrl = "ws://127.0.0.1:8001/log";
const handleSumit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  await formEl.validate((valid, fields) => {
    if (valid) {
      try {
        loginLoading.value = true;
        Launch(
          loginForm.target,
          loginForm.username,
          loginForm.password,
          loginForm.key,
          loginForm.logLevel
        );
        setTimeout(() => {
          CheckAuthed()
            .then((result: boolean) => {
              if (result) {
                loginMsg.value = "登录成功！";
                loginFlag.value = true;
                initSocket(wsUrl);
              } else {
                ElMessage.error(`登录失败：请检查输入信息!`);
              }
              loginLoading.value = false;
            })
            .catch((err: any) => {
              ElMessage.error(`登录失败：${err.message}`);
              loginLoading.value = false;
            });
        }, 1000);
      } catch (err: any) {
        ElMessage.error(`登录失败：${err.message}`);
      }
    }
  });
};

const handleLogout = () => {
  Logout();
  cancelWS();
  loginFlag.value = false;
  ElMessage.success("退出登录成功！");
};

// websocket 连接
let socket: WebSocket | undefined;
const initSocket = (url: any) => {
  if (!socket) {
    socket = new WebSocket(url);
    // 连接失败
    socket.onerror = onerrorWS;
    // 接收消息
    socket.onmessage = onmessageWS;
  }
};
const onerrorWS = () => {
  ElMessage.error(`websocket连接失败!`);
};
const onmessageWS = (e: any) => {
  noticeMsg.value.splice(1, 0, e.data);
};
const cancelWS = () => {
  socket?.close();
  socket = undefined;
  noticeMsg.value = [
    "欢迎使用Virtuallan客户端，让您的系统接入更安全，体验更好！",
  ];
};
</script>

<style lang="less" scoped>
.login-cls {
  display: flex;
  flex-direction: column;
  flex: 1;
  .login-top {
    display: flex;
    align-items: center;
    flex: 0 0 auto;
    padding: 0 20px;
    height: 50px;
    background-color: var(--vpn-top-color);
    color: var(--vpn-title-color);
    .system-title {
      font-size: 18px;
    }
  }
  .login-home {
    display: flex;
    flex: 1;
    padding: 20px 60px;
    .home-left {
      flex: 1 0 auto;
      margin-right: 40px;
      width: 100px;
      text-align: left;
      .notice-info {
        line-height: 24px;
      }
      .home-notice {
        padding: 20px;
        width: 100%;
        border-radius: 8px;
        background: #79777731;
      }
    }
    .home-right {
      flex: 2 0 auto;
      width: 100px;
      color: #121212;
      margin-left: 40px;
      .login-box {
        border-radius: 8px;
        background: #f6f7f8;
        display: flex;
        flex-direction: column;
        align-items: start;
        justify-content: start;
        padding: 20px 40px;
        p {
          width: 100%;
          text-align: left;
        }
        .box-cls {
          flex: 1;
          width: 100%;
          margin-top: 10px;
        }
      }
    }
  }
  .notice-title {
    font-size: 16px;
    font-weight: 500;
  }
  .login-title {
    font-size: 18px;
    font-weight: bold;
  }
  .login-button {
    width: 100%;
    margin-top: 10px;
  }
}
</style>

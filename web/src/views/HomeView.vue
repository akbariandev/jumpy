<script>
import axios from 'axios';
export default {

    data() {
        return {
            phoneNumber: "",
            password: "",
            snackbar: false,
            snackbarContent: null,
            snackbarTimeout: 3000,
            snackbarColor: "blue",
        }
    },
    methods: {
        runApplication(event) {
            //var loginData = {phone: this.phoneNumber, password: this.password, pageCode: localStorage.getItem('page_code')}
            axios({ method: "GET", url: import.meta.env.SERVERURL+"/run?nodes=2&group=krypton", withCredentials: false, headers: {"content-type": "application/json" } }).then(result => {

            }).catch( error => {

            });
        },
    }
}
</script>

<template>
    <button id="rubAppBtn" v-on:click="runApplication">Run</button>
    <div id="console">
    </div>
</template>
<style scoped>
#rubAppBtn{
    position: absolute;
    top: 10px;
    left: 10px;
}
#console{
    background: #000;
    color: #c40;
    width: 100%;
    height: 30%;
    position: absolute;
    bottom: 0;
    left: 0;
    padding: 20px;
}
</style>

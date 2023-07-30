let serverURL = document.currentScript.dataset.server
//ポートを変えた同一ホスト名のURLには対応できないかも?(バグるかも)
// let currentURL = location.protocol + "://" + location.host + location.pathname
let currentURL = "https://powerfulfamily.net/"
document.head.insertAdjacentHTML("beforeend",`
<style>
  .like-container{
    display: block;
    margin-left: auto;
    margin-right: auto;
    text-align: center;
  }
</style>
`)
document.currentScript.insertAdjacentHTML("afterend",`
<div class="like-container">
  <div class="like-container-like-button" onclick="Like.increment(1)">👍</div>
  <p id="like-container-like-counter"></p>
</div>
`)
class Like {
    like
    error
    constructor() {
        this.like = fetch(serverURL+"?url="+currentURL,{method:"GET"}).then((response)=> {
            if (response.ok) {
                return response.json().then((data) => {
                    return data.like
                })
            } else {
                throw new Error("Response not 200")
            }

        })
        this.write()
    }
    write(){
        if (this.like === undefined) {
            document.getElementById("like-container-like-counter").innerHTML = "ERROR"
        } else {
            document.getElementById("like-container-like-counter").innerHTML = this.like
        }
    }
    static increment(inc){
        this.like =+ inc
    }
}
    let like = new Like()
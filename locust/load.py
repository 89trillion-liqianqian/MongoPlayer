from locust import HttpUser, between, task


class WebsiteUser(HttpUser):
    wait_time = between(5, 15)



    @task
    def keyword(self):
        url = '/playerLogin'
        data = {"uniCode":"uniCode001"}
        self.client.post(url=url,data=data)
        #self.client.get("/getGift?code=jG7a4lo8")
        url = '/checkCode'
        data = {"code":"90KKHauh","uid":"8a601a2f-e101-437a-baa4-af37783c38f7"}
        self.client.post(url=url,data=data)
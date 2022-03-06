package main

import (
	"fmt"

	"github.com/smartwalle/alipay/v3"
)

func main() {
	privateKey := "MIIEpAIBAAKCAQEAiGK6ebfi8kvBR4P2ZeJwC9tUxTkXCaA2TaCBMQC6b47TpXmQmO3LqpL0zZskmyvAIHV7m5xXJxdjf2IB3F5sM+8FrSVZNNMIfKgZZQXV72VvrlOcglJqVgO8ur41wjwzbRTuLkUBBhtGTiww6GT6g8gUQHSRO08FM5v1Bpt1ODIHUQ4aoPfQtJtYrJfWfJXFPQsSUWRVB1p4BIkRQTC2QpGoynizI2R+RsPAXg8SvYv3a9kuS9GeK83jtazVDP3uCCn3KoX/Zns/IzCqzt5oaFDeECoN9fperm7dNFtDYKaBHdTpWYfUPXvAN8osuQYpFKBnlPA6h00/pTnftigfQwIDAQABAoIBAH/EnhVxY7bBqWwqACcqnvU8oImzFhX6dWT7vIv/tqjUEQqwe/oIrAg8bAJCGWUU84R1CzuYlz0T+C5Zzox4GL3+aj7vAkcufIDrqf6uHOGNEbcqtTCRrPSvw5rentY8N6E1Pn0yP+dMneyWEatlDhtu0SohCf2Z5bGXpVBvSr09Ish44PyskKrg/CvCvj3j4cZZYJK19rZATMU295WKMc3jvw91K52Hu2UIY5hazL9v56cLFsgriZRSy7Ev/99DJyhTwq49ksZUUSuHG2Y8Stsgq7OM9hru8FzmDQKN/MyV0AtOamDMXO96F8EHDL2aXRw3unBc8RIjFML+UMxJ6qECgYEAwCcrr1Rgr0v7mJ4Ja/H70ZfHKXeqqsVK3a9Ldz4n3JR6xIeDy61EhWIZp0EVH+IRkXZT6qqu/5jXOwVrQBHD/7XZnTQxGoCxcu03Yaug2YOl+iftKNfvQaITjs5Y6FpYZGvPxti7WSaxQtcA0sRJ1Vg94eZn8OtXv2i+KyyL5/0CgYEAtbPmv960h4AzYGbkCStGNbeIKbaW6udIjDdv4mDT4bxGxaBSkzkiSWXOzhs7V0UZ3cQS3+k+bptieCulVRJuKIPv776XNDz4ndyTRGDvmQlHcMpdHHUZ9dqtfZAkv+LaKOZDjkiY2nD57UuTifAJOIn6hR1DpimkpcmwLVAsqD8CgYEAmx9EMf/RKdMSYsu6WW152GNKQhy8J66sWLjKGJKSBY3MalnoOQZA2dkvUonE7v9HJYI8DqcKLYeKwbgHNCrjasy4yCM5POcF2fzNB8lRSifwVzniSGXCXd7lIRVOSw2cbD1o+GNBI6CL26TMoloaLORW2MZzxNeI+Bor45jLvVUCgYEAs3FxRrdXriGrm17BgVSdR9tyu085B89VVDRDaFubpGjds7o7Em3wMHA8pks6dVsmyl4jDcI5B96ohmkEJFnJNHXn9OpSRSKZnL0DKxpYRNhnFzqibcIv2x8VCtXZlS8hqBaPTOrhGYlNKU3j7OuDD7UkFWXrMyQZGClwta9iCt0CgYAIOUwXOI9sTpxd5zvTIIyrWcsAq17HcFEZbnkar47rFogZ2fUsqzn9z8yIStUycZZqhwJYn+rT3dXfbDmsgwbrv6n9rkTNGR8AS/iGbtWOsjZ7P8br0+Fu9c5/h/etCxYuet9WXBLW0QdUjyo1qDl5SV17AYktKf0ih5Otd59+2Q=="
	client, err := alipay.New("2021000119628913", privateKey, false)
	if err != nil {
		panic(err)
		return
	}

	publicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsF33Gybk1duYmSbvj7bZp37vyYRTSv2zjB7YjCri8GLRVAU+OHyfy/tVwwIAL/5SzKWr1SHFKJBy8AwOHWLT3gD6Qt/LqmY2gtIZTjeupfQ8lc21HUGeBITONuLHkwMsGVztXb4GSDl0bi1AS6wATQtKmMRoWJdpBO4Eqhne6KokoSk6QB35W760zpjzwRMBD7y84+koadgwE8ySb/JEjn4JSfNLywyo4CPzvPEK5gL/LP7hivaSge2eJ+4SJ/NRvENEnA0Mbu7+UN750jsgLhBOLnkNE6UJDRzw6K+apxidb9wVeJGiMHlaAXX/xhm1PsIy//pgXT1hH6sJvkat7wIDAQAB"
	err = client.LoadAliPayPublicKey(publicKey)
	if err != nil {
		panic(err)
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = "http://xxx"
	p.ReturnURL = "http://xxx"
	p.Subject = "标题"
	p.OutTradeNo = "传递一个唯一单号"
	p.TotalAmount = "10.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	payURL, err := client.TradePagePay(p)
	if err != nil {
		return
	}
	fmt.Println(payURL.String())
}

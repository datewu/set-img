custom set-img
2020年10月12日 星期一 14时52分21秒 CST
kku --context hk  create secret generic set-img --from-file=private_key_for_sign.pem=private_key_for_sign.pem
kku --context hk  apply -f sa.yaml
kku --context hk  apply -f deployment.yaml  
kku --context hk expose deploy set-img --port 80 --target-port 8080

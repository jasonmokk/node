# Messages

## MsgUpdatePolicies

UpdatePolicies updates policies

```proto
message MsgUpdatePolicies {
	string creator = 1;
	Policies policies = 2;
}
```

## MsgUpdateChainInfo

UpdateChainInfo updates the chain info structure that adds new static chain info or overwrite existing chain info
on the hard-coded chain info

```proto
message MsgUpdateChainInfo {
	string creator = 1;
	ChainInfo chain_info = 2;
}
```


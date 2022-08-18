# meka-dev/tendermint

This fork of [tendermint/tendermint](https://github.com/tendermint/tendermint)
has been patched to supoort the [Mekatek builder API](https://api.mekatek.xyz).
Each supported network has a tracking branch, which participating validators
should use when building their full node.

<table style="text-align: left;">
  <tr>
    <th>Network</th>
    <th>Network version</th>
    <th>Tendermint version</th>
    <th>Patch version</th>
    <th>Diff</th>
  </tr>
  <tr>
    <td>Osmosis</td>
    <td><a href="https://github.com/osmosis-labs/osmosis/tree/v11.0.1">v11.0.1</a></td>
    <td><a href="https://github.com/osmosis-labs/osmosis/blob/v11.0.1/go.mod#L28">v0.34.9</td>
    <td><a href="https://github.com/meka-dev/tendermint/tree/v0.34.x-meka">v0.34.x-meka</a></td>
    <td><a href="https://github.com/meka-dev/tendermint/compare/v0.34.x...v0.34.x-meka">diff</a></td>
  </tr>
</table>

Here is an example of how to build Osmosis with a patched Tendermint.

```shell
git clone https://github.com/osmosis-labs/osmosis
cd osmosis
git checkout v11.0.1
go mod edit -replace=github.com/tendermint/tendermint=github.com/meka-dev/tendermint@v0.34.x-meka
go mod tidy
make install
```

<sub>[upstream README](/README.upstream.md)</sub>

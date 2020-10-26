package paillier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCiphertextAdd(t *testing.T) {
	c1 := "010070061bde0fc57d38e0ce9c82c088f793a93a20f830ed94a5f7e01f5af0decccf0ee7f70a711dc78b69b40f7e411f7f6ff34a5de4ad32fa5dc9d09b09dbbcf8e49bb1c9be19b215577167458b969bb40634f71dba00bdfdd12f1c878effd9e08475b8301a214bc7dfc95c2f341aa8fbd985a5da05af046336c23cf356570772530017c9a8bb96140977698d50c3258e8d31c040bb887c56d3d6e583db7a4ff2ea28aa93885312dfa50bf86a15618288b2c2701e6d75743fec9f3ca35b7e20ad902ad6985fe0b3c389b9582b3c62d9bea90230d1512a07c66a5404bf804693b6bb677a4ff038ae4a762f89f6a877d7621f7461431a7d4cadaf1687057ca4a9d1570e4ceeab32c86856d89731d50edf6f0886ccf348cc4565227f1df6b813569c49b3170e9efd722e4f4e957add9a115a2122a9f5c00780a1b8020d44b8b865d16113eadcfd15d1b755aed7c983b3d4b9390eda68d22d165a32d4055eee95adee79d7d28d2ee6d593d4eff217cd3e3d70bb9c32a6608af78afb7d5cb7ee0422500273b934209b950cd07842ab7e10c0dfc790d212c9085058d39c667b7dfff0963478118451ec64af30aff518384cc47758b2aa4331c8b5f36d6d1b9c3a6ee7ba77d28ea6ca844960e24ab6384e36127745f1285f3b71cea6c5c6a41c51ca7edbfae07766a3fb7bd9ef6b8f999460741a5b44eeb0b98937b57f4582294e18fd70d402dee2e557fa7df5cfb13240458f77b62a32169940ab5a1fce95b52583ecb13e21c22749f7620caaa6fd998e4f9bd0e617c925b83bed5ed2215ec8a404357a782b70d2c6ece6fa7a4fb30d18a01b81b2f3a8e38cc12f32becfdf195e02879b710f2feda5b0c2e86ba6df2a0fa25c60cef58224ee466787034cc5d0872a28cf8d6bad106125660b0b4bfe86eacb829dbb073420ff3af3ac76e5718c199fcf53258397a9e4fed0d98f360d16248693aced1ecbc5ca5211f48df011d14d372ee07f9b7c24bc219f3cdc403ebeb5b452643ddd1f4276779fab76a19c8b195a1213a235f7facf67a26b8a7804cb326f58200d756d62d40edb650498c198a524c091f0"

	c2 := "010070061bde0fc57d38e0ce9c82c088f793a93a20f830ed94a5f7e01f5af0decccf0ee7f70a711dc78b69b40f7e411f7f6ff34a5de4ad32fa5dc9d09b09dbbcf8e49bb1c9be19b215577167458b969bb40634f71dba00bdfdd12f1c878effd9e08475b8301a214bc7dfc95c2f341aa8fbd985a5da05af046336c23cf356570772530017c9a8bb96140977698d50c3258e8d31c040bb887c56d3d6e583db7a4ff2ea28aa93885312dfa50bf86a15618288b2c2701e6d75743fec9f3ca35b7e20ad902ad6985fe0b3c389b9582b3c62d9bea90230d1512a07c66a5404bf804693b6bb677a4ff038ae4a762f89f6a877d7621f7461431a7d4cadaf1687057ca4a9d157195acafe73f5eed77b314dd640999ee20de461600ce3d3620cef84abc256774a514db5b65721cd93d1601e0fa8f53630dbe8567493018588e9619783e1526eb9bc56cc7b37c99025a9bdc10465bfb01605232292e50ac3b56b2ffefe48f0615651649d1db4ebff57d955e434f0c564b2a6842d3af04be1bfc5e968abe390f15bc89b9a4ed7a00b82beea07cf8e03c82bdd6a40b6c6334e51a432024caf9247db25058eb0cbce567ed49c56730058480e5854f3888e988053c403045239098ba6977aa0adb7f6dec785d31bacee8276a87d7a8f2114bb6662a117e3e37846c5b713e599801fbae5185cf21360760721376efe0736b7f4fa049d0c9f3086324be3a0f44d6cae5450d87d79a7bcaa0a2aa259b738ed85a59fead48d985597f6d77022ecdf37685df4182267d4fbdb24747f208a9e1b3126251b320066d8499ab5298395dab0340a78325b39d383d34a31a27f223114a9774498a314e3ba3cf4b87150e1f959f29556e67ae136a6290922e8cd3890c3b93bb68f38990ac3a4700cbe075795f7603755a28f3226e31c29a2eab2bff2423dc6177a6d8eaa950fbe2320b9bee89b8bbc3149b7ee177224e59726f28bcb796ccb09ac6b3feb3256ca033eb177aeea997be798f0e1b63b284dc8361ff870de717263cb86838a0c6cd5a82e12d162bd04f9c85267d4abd3a6828639fd0584a023388cc220d404d4646d0318"

	c3 := "010070061bde0fc57d38e0ce9c82c088f793a93a20f830ed94a5f7e01f5af0decccf0ee7f70a711dc78b69b40f7e411f7f6ff34a5de4ad32fa5dc9d09b09dbbcf8e49bb1c9be19b215577167458b969bb40634f71dba00bdfdd12f1c878effd9e08475b8301a214bc7dfc95c2f341aa8fbd985a5da05af046336c23cf356570772530017c9a8bb96140977698d50c3258e8d31c040bb887c56d3d6e583db7a4ff2ea28aa93885312dfa50bf86a15618288b2c2701e6d75743fec9f3ca35b7e20ad902ad6985fe0b3c389b9582b3c62d9bea90230d1512a07c66a5404bf804693b6bb677a4ff038ae4a762f89f6a877d7621f7461431a7d4cadaf1687057ca4a9d1571d74c0a358b73242b65cf84139071503ffa45c824d84cb0336510de9a1ad1a7c56142efa23f8b69b4094c940adb42df385fc11bf2704a2608b1a65c5d17f2a051114832ab141290c16201b59eae28639ece329eb5b8003b4f1025e57066115d7d378581268359e0ca94117451ca3aa6b0d597cf8f621b66bf20fc488b18848aa47e31d5939d59b8058def7df416f6734b3821bacd1fc8f4f1d422fc2c17c9bd903e80fabea7eb0f6237f771d0be65941cfca8ed6ec0e412bc940fc09cb0033abe4e668fee7f07a9cba77da04815747271eb32a5a9278ca0f57d7e05cffa6bdb5902984586259fab80be0460e9af3391f2a3212a1e4f350e060f988781e6ab2254e76c972e5611151fa501607d1d342667253d21e9ca113f2b87d823a55e70d4f5899e0ba8121477dded3d3b858438c4d0486a0c9846eb05b737a734ba37cdfdaaca8a4a38f2bd1c1ed0b4583b96c2fd2d1675991b4a151470811b4f4bfbad46ddbaa0e1fa28a81a3801e8c211137db822fb5a07344dc719b13e50562bd1e3a391a7747709db2445543f7412fd9c3011141d9ad4f05182593e6e75033764395fa8be80039cb47ed5d750f62836312e7b23152f69a65e41a9782d841b3c48719403e98bd8e72c824083275a29a53d681dc25b7df6c60cb0c4c74feab1d7b1edc76c09a1fe4da1acf9979d3c200233b082f66710725867ef6418e866cd93d26c7fd"

	data, err := CiphertextAdd(c1, c2)
	assert.Nil(t, err)
	assert.Equal(t, c3, data)
}

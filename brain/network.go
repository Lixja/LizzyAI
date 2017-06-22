package brain

import "math/rand"
import "fmt"
import "time"
import "strconv"

type Network struct {
	Deep [][]Neuron
}

func (n *Network) Train(tvalue float64, InputV, WantedResult []float64, iterations int, fast bool) bool {
	write("Training Brain")
	if fast {
		write("FastTraining")
	} else {
		write("SlowTraining")
	}
	var backup *Network
	var backupw []float64
	var dp int
	var dlp int
	for c := 0; c < iterations; c++ {
		Results := n.Think(InputV)
		difference := 0.0
		pon := 0
		for i, res := range WantedResult {
			difference += positive(res - Results[i])
			if res > Results[i] {
				pon++
			} else if res < Results[i] {
				pon--
			}
		}
		if pon < 0 {
			tvalue = -tvalue
		} else if pon == 0 {
			return true
		}
		if fast {
			backup = n.Backup()
			for i := 0; i < len(n.Deep); i++ {
				for i2 := 0; i2 < len(n.Deep[i]); i2++ {
					n.Deep[i][i2].modify(tvalue)
				}
			}
		} else {
			dp = rand.Intn(len(n.Deep))
			dlp = rand.Intn(len(n.Deep[dp]))
			backupw = make([]float64, len(n.Deep[dp][dlp].Weight))
			copy(backupw, n.Deep[dp][dlp].Weight)
			n.Deep[dp][dlp].modify(tvalue)
		}
		Results = nil
		Results = n.Think(InputV)
		ndifference := 0.0
		for i, res := range WantedResult {
			ndifference += positive(res - Results[i])
		}
		if ndifference > difference {
			if fast {
				for i := 0; i < len(n.Deep); i++ {
					for i2 := 0; i2 < len(n.Deep[i]); i2++ {
						n.Deep[i][i2].Weight = backup.Deep[i][i2].Weight
					}
				}

			} else {
				copy(n.Deep[dp][dlp].Weight, backupw)
			}
		}
		for i := 0.0; i < 100.0; i += 10 {
			if (float64(c) / (float64(iterations) / 100)) == i {
				write("Training completed: " + strconv.FormatFloat(i, 'f', -1, 64) + "%")
			}
		}
	}
	return false
}

func GenerateBrain(ainputs, aoutputs, alayers, mnpl int) *Network {
	write("Generating Brain")
	write("Creating Neurons")
	var net Network
	rand.Seed(time.Now().UTC().UnixNano())
	var layer []Neuron
	for i := 0; i < ainputs; i++ {
		layer = append(layer, Neuron{})
	}
	net.Deep = append(net.Deep, layer)
	for i := 0; i < alayers; i++ {
		layer = nil
		for i2 := rand.Intn(mnpl - 1); i2 < mnpl; i2++ {
			layer = append(layer, Neuron{})
		}
		net.Deep = append(net.Deep, layer)
	}
	layer = nil
	for i := 0; i < aoutputs; i++ {
		layer = append(layer, Neuron{})
	}
	net.Deep = append(net.Deep, layer)
	write("Created Neurons")
	write("Connecting Neurons")
	for i := 0; i < len(net.Deep); i++ {
		for i2 := 0; i2 < len(net.Deep[i]); i2++ {
			if i+1 < len(net.Deep) {
				for i3 := 0; i3 < len(net.Deep[i+1]); i3++ {
					net.Deep[i][i2].Connect(&net.Deep[i+1][i3])
				}
			}
		}
		write(strconv.Itoa(i+1) + " of " + strconv.Itoa(len(net.Deep)) + " layers ready.")
	}
	write("Connected Neurons")
	return &net
}

func (n *Network) Think(val []float64) []float64 {
	var output []float64
	for i := 0; i < len(n.Deep[0]); i++ {
		n.Deep[0][i].DoTheThing(val[i], &output)
	}
	for len(output) < len(n.Deep[len(n.Deep)-1]) {
		time.Sleep(100 * time.Microsecond)
	}
	return output
}

func write(msg string) {
	fmt.Println("Brain: " + msg)
}

func positive(i float64) float64 {
	if i >= 0 {
		return i
	} else {
		return -i
	}
}

func (n *Network) Backup() *Network {
	var backup Network
	backup.Deep = make([][]Neuron, len(n.Deep))
	for i := 0; i < len(n.Deep); i++ {
		backup.Deep[i] = make([]Neuron, len(n.Deep[i]))
		copy(backup.Deep[i], n.Deep[i])
		for i2 := 0; i2 < len(n.Deep[i]); i2++ {
			backup.Deep[i][i2].Weight = make([]float64, len(n.Deep[i][i2].Weight))
			copy(backup.Deep[i][i2].Weight, n.Deep[i][i2].Weight)
		}
	}
	return &backup
}

import numpy as np
import bpy

def draw_voxel_model(V, N, M, K, group_name='VM'):
  g = bpy.data.groups.new(group_name)
  mat_dict = dict()
  Nh = N/2
  Mh = M/2
  Kh = K/2
  for i in range(N):
    for j in range(M):
      for k in range(K):
        p = V[i,j,k]
        if p > 0:
          mat_name = 'm'+str(p)[2:]
          mat = bpy.data.materials.get(mat_name)
          if mat is None:  
            mat = bpy.data.materials.new(mat_name)
            mat.diffuse_color = (0.5,0.5,0.5)
            mat.alpha = p
            mat.use_transparency = True
          bpy.ops.mesh.primitive_cube_add(location=(i+1/2-Nh, j+1/2-Mh, k+1/2-Kh))
          v = bpy.context.active_object
          v.dimensions = (1,1,1)
          v.active_material = mat
          v.show_transparent = True
          g.objects.link(v)
  return g

